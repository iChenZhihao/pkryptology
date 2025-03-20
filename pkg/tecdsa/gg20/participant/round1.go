//
// Copyright Coinbase, Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0
//

package participant

import (
	"fmt"
	"math/big"

	"github.com/coinbase/kryptology/internal"
	"github.com/coinbase/kryptology/pkg/core"
	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/tecdsa/gg20/proof"
)

// Round1Bcast contains values to be broadcast to all players after the completion of singing round 1
type Round1Bcast struct {
	Identifier uint32
	C          core.Commitment
	Ctxt       *big.Int
	Proof      *proof.Range1Proof
}

type Round1P2PSend = proof.Range1Proof

// SignRound1 performs round 1 signing operation
// Trusted Dealer Mode: see [spec] fig 7: SignRound1
// DKG Mode: see [spec] fig 8: SignRound1
// NOTE: Pseudocode shows N~, h1, h2, the curve's g, q, and signer's public key as inputs
// Since `signer` already knows the paillier secret and public keys, this input is not necessary here
// `participant.PrepareToSign` receives the other inputs and stores them as state variables.
// ECDSA
func (signer *Signer) SignRound1() (*Round1Bcast, map[uint32]*Round1P2PSend, error) {
	if signer == nil || signer.Curve == nil {
		return nil, nil, internal.ErrNilArguments
	}

	if err := signer.verifyStateMap(1, nil); err != nil {
		return nil, nil, err
	}

	pk := &signer.sk.PublicKey

	// 1. k_i \getsr Z_q
	k, err := core.Rand(signer.Curve.Params().N)
	if err != nil {
		return nil, nil, err
	}

	// 2. \gamma_i \getsr Z_q
	gamma, err := core.Rand(signer.Curve.Params().N)
	if err != nil {
		return nil, nil, err
	}

	// 3. \Gamma_i = g^{\gamma_i} in \G
	Gamma, err := curves.NewScalarBaseMult(signer.Curve, gamma)
	if err != nil {
		return nil, nil, err
	}

	// 4. C_i, D_i = Commit(\Gamma_i) 计算承诺
	Ci, Di, err := core.Commit(Gamma.Bytes())
	if err != nil {
		return nil, nil, err
	}

	// 5. c_i, r_i = PaillierEncryptAndReturnRandomness(pk_i, k_i)
	ctxt, r, err := pk.Encrypt(k) // ctxt:paillier加密后的密文 r:该次加密生成的随机数
	if err != nil {
		return nil, nil, err
	}

	pp := proof.Proof1Params{
		Curve: signer.Curve, //曲线
		Pk:    pk,           //该signer的paillier公钥
		A:     k,            //该signer在本次加密选取的随机数
		C:     ctxt,         //ctxt:paillier加密后的密文
		R:     r,            // 该次加密生成的随机数
	}
	bcast := Round1Bcast{
		Identifier: signer.id, //该signer的id
		C:          Ci,        //对于Gamma_i的承诺值
		Ctxt:       ctxt,      //ctxt:paillier加密后的密文
	}
	p2p := make(map[uint32]*Round1P2PSend)

	if signer.state.keyGenType.IsTrustedDealer() {
		// 使用中心化的密钥生成
		pp.DealerParams = signer.state.keyGenType.GetProofParams(1)
		// 6. TrustedDealer - \pi_i^{\Range1} = MtAProveRange1(g,q,pk_i,N~,h_1,h_2,k_i,c_i,r_i)
		bcast.Proof, err = pp.Prove()
		if err != nil {
			return nil, nil, err
		}
	} else {
		// 使用分布式密钥生成（后续一般用它）
		// 6. (figure 8.)DKG - for j = [1,...,t+1]
		for id := range signer.state.cosigners {
			// 获取当前节点要发给其他节点的Prove证明
			// 7. DKG if i == j, continue
			if signer.id == id {
				// 为自己
				continue
			}
			pp.DealerParams = signer.state.keyGenType.GetProofParams(id)
			if pp.DealerParams == nil {
				return nil, nil, fmt.Errorf("no proof params found for cosigner %d", id)
			}
			// 8. DKG \pi_ij^{\Range1} = MtAProveRange1(g,q,pk_i,N~j,h_1j,h_2j,k_i,c_i,r_i)
			pi, err := pp.Prove()
			if err != nil {
				return nil, nil, err
			}
			// 9. P2PSend
			p2p[id] = pi
		}
	}

	// 8. Stored locally (k_i, \gamma_i, D_i, c_i, r_i)
	signer.Round = 2
	signer.state.ki = k
	signer.state.gammai = gamma
	signer.state.Gammai = Gamma
	signer.state.Di = Di
	signer.state.ci = ctxt
	signer.state.ri = r

	// (figure 7) 7. Broadcast (C_i, c_i, \pi^{Range1}_i)
	// (figure 8) 9. P2PSend(\pi^{Range1}_ij)
	// (figure 8) 10. Broadcast (C_i, c_i)
	return &bcast, p2p, nil
}
