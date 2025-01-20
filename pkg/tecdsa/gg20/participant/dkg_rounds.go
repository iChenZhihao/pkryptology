package participant

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"github.com/coinbase/kryptology/pkg/core"
	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/tecdsa/gg20/dealer"
	"github.com/stretchr/testify/require"
	"math/big"

	//"github.com/stretchr/testify/require"
	//"math/big"
	"testing"

	//"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/paillier"
	v1 "github.com/coinbase/kryptology/pkg/sharing/v1"
	//"github.com/stretchr/testify/require"
	//"math/big"
	//"testing"
	"time"
)

var (
	ecdsaVerifier = func(verKey *curves.EcPoint, hash []byte, sig *curves.EcdsaSignature) bool {
		pk := &ecdsa.PublicKey{
			Curve: verKey.Curve,
			X:     verKey.X,
			Y:     verKey.Y,
		}
		return ecdsa.Verify(pk, hash, sig.R, sig.S)
	}
)

func DkgFullRoundsWorks(t *testing.T, curve elliptic.Curve, total, threshold int) (map[uint32]*DkgParticipant, map[uint32]*DkgResult, error) {
	now := time.Now()
	fmt.Printf("start time: %v\n", now)
	var err error
	now = time.Now()
	fmt.Printf("before initiate time: %v\n", now)

	// Initiate 3 parties for DKG
	dkgParticipants := make(map[uint32]*DkgParticipant, total)
	for i := 1; i <= total; i++ {
		dkgParticipants[uint32(i)] = &DkgParticipant{
			Curve: curve,
			id:    uint32(i),
			Round: 1,
			state: &dkgstate{
				Threshold: uint32(threshold),
				Limit:     uint32(total),
			},
		}
	}

	now = time.Now()
	fmt.Printf("before dkg round1 time: %v\n", now)

	// Run Dkg Round 1
	dkgR1Out := make(map[uint32]*DkgRound1Bcast, total)
	for i := 1; i <= total; i++ {
		dkgR1Out[uint32(i)], err = dkgParticipants[uint32(i)].DkgRound1(uint32(threshold), uint32(total))
		if err != nil {
			fmt.Printf("error: %v\n", err.Error())
			return nil, nil, err
		}
	}

	now = time.Now()
	fmt.Printf("before dkg round2 time: %v\n", now)

	// Run Dkg Round 2
	dkgR2Bcast := make(map[uint32]*DkgRound2Bcast, total)
	dkgR2P2PSend := make(map[uint32]map[uint32]*DkgRound2P2PSend, total)
	for i := 1; i <= total; i++ {
		dkgR2Bcast[uint32(i)], dkgR2P2PSend[uint32(i)], err = dkgParticipants[uint32(i)].DkgRound2(dkgR1Out)
		//require.NoError(t, err)
		if err != nil {
			fmt.Printf("error: %v\n", err.Error())
			return nil, nil, err
		}
	}
	now = time.Now()
	fmt.Printf("before dkg round3 time: %v\n", now)

	// Run Dkg Round 3
	decommitments := make(map[uint32]*core.Witness, total)
	dkgR3Out := make(map[uint32]paillier.PsfProof, total)
	for i := 1; i <= total; i++ {
		decommitments[uint32(i)] = dkgR2Bcast[uint32(i)].Di
	}
	//decommitments[1] = dkgR2Bcast[1].Di
	//decommitments[2] = dkgR2Bcast[2].Di
	//decommitments[3] = dkgR2Bcast[3].Di

	for i := 1; i <= total; i++ {
		shamirMap := make(map[uint32]*v1.ShamirShare, total)
		for j := 1; j <= total; j++ {
			shamirMap[uint32(j)] = dkgParticipants[uint32(j)].state.X[uint32(i-1)]
		}
		dkgR3Out[uint32(i)], err = dkgParticipants[uint32(i)].DkgRound3(decommitments, shamirMap)
		if err != nil {
			fmt.Printf("error: %v\n", err.Error())
			return nil, nil, err
		}
	}
	//dkgR3Out[1], err = dkgParticipants[1].DkgRound3(decommitments, map[uint32]*v1.ShamirShare{
	//	1: dkgParticipants[1].state.X[0],
	//	2: dkgParticipants[2].state.X[0],
	//	3: dkgParticipants[3].state.X[0],
	//})
	//require.NoError(t, err)

	//dkgR3Out[2], err = dkgParticipants[2].DkgRound3(decommitments, map[uint32]*v1.ShamirShare{
	//	1: dkgParticipants[1].state.X[1],
	//	2: dkgParticipants[2].state.X[1],
	//	3: dkgParticipants[3].state.X[1],
	//})
	//require.NoError(t, err)
	//
	//dkgR3Out[3], err = dkgParticipants[3].DkgRound3(decommitments, map[uint32]*v1.ShamirShare{
	//	1: dkgParticipants[1].state.X[2],
	//	2: dkgParticipants[2].state.X[2],
	//	3: dkgParticipants[3].state.X[2],
	//})
	//require.NoError(t, err)

	now = time.Now()
	fmt.Printf("before dkg round4 time: %v\n", now)

	// Run Dkg Round 4
	dkgR4Out := make(map[uint32]*DkgResult, total)
	for i := 1; i <= total; i++ {
		dkgR4Out[uint32(i)], err = dkgParticipants[uint32(i)].DkgRound4(dkgR3Out)
		//require.NoError(t, err)
	}

	now = time.Now()
	fmt.Printf("finish dkg rounds time: %v\n", now)

	//// Check that the shares result in valid secret key and public key
	field := curves.NewField(curve.Params().N)
	//
	shamir, _ := v1.NewShamir(threshold, total, field)
	share1 := v1.NewShamirShare(1, dkgR4Out[1].SigningKeyShare.Bytes(), field)
	share2 := v1.NewShamirShare(2, dkgR4Out[2].SigningKeyShare.Bytes(), field)
	share3 := v1.NewShamirShare(3, dkgR4Out[3].SigningKeyShare.Bytes(), field)
	//share4 := v1.NewShamirShare(4, dkgR4Out[4].SigningKeyShare.Bytes(), field)
	//share5 := v1.NewShamirShare(5, dkgR4Out[5].SigningKeyShare.Bytes(), field)

	//secret12, err := shamir.Combine(share1, share2)
	//require.NoError(t, err)
	//secret23, err := shamir.Combine(share2, share3)
	//require.NoError(t, err)

	secret123, err := shamir.Combine(share1, share2, share3)
	//require.NoError(t, err)
	//secret134, err := shamir.Combine(share1, share3, share4)
	//require.NoError(t, err)
	//secret235, err := shamir.Combine(share2, share3, share5)
	//require.NoError(t, err)

	//require.Equal(t, secret12, secret23)

	// Check the relationship of verification key and signing key is valid
	pk, err := curves.NewScalarBaseMult(curve, new(big.Int).SetBytes(secret123))
	require.NoError(t, err)
	require.True(t, dkgParticipants[1].state.Y.Equals(pk))
	//
	//// Check every participant has the same verification key
	//require.Equal(t, dkgParticipants[1].state.Y, dkgParticipants[2].state.Y)
	//require.Equal(t, dkgParticipants[1].state.Y, dkgParticipants[3].state.Y)
	//
	//// Testing validity of paillier public key and secret key
	//// Check every participant receives equal paillier public keys from other parties
	//require.Equal(t, dkgParticipants[1].state.otherParticipantData[2].PublicKey, dkgParticipants[3].state.otherParticipantData[2].PublicKey)
	//require.Equal(t, dkgParticipants[1].state.otherParticipantData[3].PublicKey, dkgParticipants[2].state.otherParticipantData[3].PublicKey)
	//require.Equal(t, dkgParticipants[2].state.otherParticipantData[1].PublicKey, dkgParticipants[3].state.otherParticipantData[1].PublicKey)
	//
	//// Testing validity of paillier keys of participant 1
	//pk1 := dkgParticipants[2].state.otherParticipantData[1].PublicKey
	//sk1 := dkgParticipants[1].state.Sk
	//msg1, _ := core.Rand(pk1.N)
	//c1, r1, err := pk1.Encrypt(msg1)
	//require.NoError(t, err)
	//require.NotNil(t, c1, r1)
	//m1, err := sk1.Decrypt(c1)
	//require.Equal(t, m1, msg1)
	//require.NoError(t, err)
	//
	//// Testing validity of paillier keys of participant 2
	//pk2 := dkgParticipants[1].state.otherParticipantData[2].PublicKey
	//sk2 := dkgParticipants[2].state.Sk
	//msg, _ := core.Rand(pk2.N)
	//c, r, err := pk2.Encrypt(msg)
	//require.NoError(t, err)
	//require.NotNil(t, c, r)
	//m, err := sk2.Decrypt(c)
	//require.Equal(t, m, msg)
	//require.NoError(t, err)
	//
	//// Testing validity of paillier keys of participant 3
	//pk3 := dkgParticipants[1].state.otherParticipantData[3].PublicKey
	//sk3 := dkgParticipants[3].state.Sk
	//msg3, _ := core.Rand(pk3.N)
	//c3, r3, err := pk3.Encrypt(msg3)
	//require.NoError(t, err)
	//require.NotNil(t, c3, r3)
	//m3, err := sk3.Decrypt(c3)
	//require.Equal(t, m3, msg3)
	//require.NoError(t, err)
	//
	//// Checking public shares are equal
	//require.Equal(t, dkgParticipants[1].state.PublicShares, dkgParticipants[2].state.PublicShares)
	//require.Equal(t, dkgParticipants[1].state.PublicShares, dkgParticipants[3].state.PublicShares)
	//
	//// Checking proof params are equal
	//require.Equal(t, dkgR4Out[1].ParticipantData[2].ProofParams, dkgR4Out[3].ParticipantData[2].ProofParams)
	//require.Equal(t, dkgR4Out[1].ParticipantData[3].ProofParams, dkgR4Out[2].ParticipantData[3].ProofParams)
	//require.Equal(t, dkgR4Out[2].ParticipantData[1].ProofParams, dkgR4Out[3].ParticipantData[1].ProofParams)

	return dkgParticipants, dkgR4Out, err
}

func ConvertToSigners(dkgParticipants map[uint32]*DkgParticipant, dkgResults map[uint32]*DkgResult, threshold uint) (*curves.EcPoint, map[uint32]*Signer) {
	var pk *curves.EcPoint
	//var err error
	signers := make(map[uint32]*Signer, threshold)
	//var cosigners []uint32
	cosigners := []uint32{1, 2, 3}
	encryptKeys := make(map[uint32]*paillier.PublicKey, threshold)
	proofParams := make(map[uint32]*dealer.ProofParams, threshold)
	publicShare := make(map[uint32]*dealer.PublicShare, threshold)
	shareMap := make(map[uint32]*dealer.Share, threshold)

	for _, id := range cosigners {
		pk = dkgParticipants[id].state.Y
		encryptKeys[id] = &dkgResults[id].EncryptionKey.PublicKey
		proofParams[id] = &dealer.ProofParams{
			N:  dkgParticipants[id].state.N,
			H1: dkgParticipants[id].state.H1,
			H2: dkgParticipants[id].state.H2,
		}

		publicSharePoint, _ := curves.NewScalarBaseMult(pk.Curve, dkgParticipants[id].state.Xii.Value.BigInt())
		shareMap[id] = &dealer.Share{
			Point:       publicSharePoint,
			ShamirShare: dkgParticipants[id].state.Xii,
		}
	}

	publicShare, _ = dealer.PreparePublicShares(shareMap)

	for _, id := range cosigners {
		p := &dealer.ParticipantData{
			Id:             id,
			EcdsaPublicKey: dkgParticipants[id].state.Y,
			EncryptKeys:    encryptKeys,
			DecryptKey:     dkgResults[id].EncryptionKey,
			KeyGenType:     &dealer.DistributedKeyGenType{ProofParams: proofParams},
			PublicShares:   publicShare,
		}

		p.SecretKeyShare = shareMap[id]

		signer, err := NewSigner(p, ecdsaVerifier, cosigners)
		if err != nil {
			fmt.Printf("New Signer Error: %v\n", err)
		}
		signers[id] = signer
		signers[id].threshold = threshold
	}

	return pk, signers
}
