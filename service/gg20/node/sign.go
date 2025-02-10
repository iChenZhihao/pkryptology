package node

import (
	"fmt"
	"github.com/coinbase/kryptology/pkg/core"
	"github.com/coinbase/kryptology/pkg/core/curves"
	ptcpt "github.com/coinbase/kryptology/pkg/tecdsa/gg20/participant"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var (
	signOnce     sync.Once
	signOperator *SignOperator
)

const SignRoundWait = 37

func GetSignOperator() *SignOperator {
	signOnce.Do(func() {
		signOperator = &SignOperator{
			signerMap: make(map[string]*SignerInfo, 2000),
		}
	})
	return signOperator
}

type SignOperator struct {
	id                uint32
	threshold         uint32
	total             uint32
	participantAddrs  []string
	otherParticipants []uint32
	signerMap         map[string]*SignerInfo
	cond              sync.RWMutex
}

type SignerInfo struct {
	id                uint32
	workId            string
	cosigner          []uint32
	participantAddrs  []string
	hashMsg           []byte
	candidateInfoChan chan *CandidateInfo
	round1RecvChan    chan SignRound1Recv
	round2RecvChan    chan SignRound2Recv
	round3RecvChan    chan SignRound3Recv
	round4RecvChan    chan SignRound4Recv
	round5RecvChan    chan SignRound5Recv
	round6RecvChan    chan SignRound6Recv
	signer            *ptcpt.Signer
}

func (s *SignOperator) SignMsg(plainMsg []byte) (*curves.EcdsaSignature, error) {
	s.cond.RLock()
	defer s.cond.RUnlock()
	if !GetDkgOperator().IsAvailable() {
		return nil, fmt.Errorf("签名节点集群暂不可用")
	}
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	signerInfo := &SignerInfo{
		id:                s.id,
		workId:            newUUID.String(),
		cosigner:          []uint32{s.id},
		candidateInfoChan: make(chan *CandidateInfo, s.total+1),
		round1RecvChan:    make(chan SignRound1Recv, s.threshold+1),
		round2RecvChan:    make(chan SignRound2Recv, s.threshold+1),
		round3RecvChan:    make(chan SignRound3Recv, s.threshold+1),
		round4RecvChan:    make(chan SignRound4Recv, s.threshold+1),
		round5RecvChan:    make(chan SignRound5Recv, s.threshold+1),
		round6RecvChan:    make(chan SignRound6Recv, s.threshold+1),
		participantAddrs:  s.participantAddrs,
	}
	s.signerMap[signerInfo.workId] = signerInfo
	s.DoAskForCosignerCandidate(signerInfo.workId)

	for {
		if len(signerInfo.cosigner) == int(s.threshold) {
			break
		}
		select {
		case recv := <-signerInfo.candidateInfoChan:
			signerInfo.cosigner = append(signerInfo.cosigner, recv.Id)
		case <-time.After(SignRoundWait * time.Second):
			return nil, fmt.Errorf("%d节点等待确定其它cosigner超时", s.id)
		}
	}
	signerInfo.signer, err = GetDkgOperator().NewSigner(signerInfo.cosigner)
	if err != nil {
		return nil, err
	}

	glog.Info("确定了cosigner：", s.signerMap[signerInfo.workId].cosigner)
	hashMsg, err := core.Hash(plainMsg, curve)
	if err != nil {
		return nil, err
	}
	signerInfo.hashMsg = hashMsg.Bytes()
	var wg sync.WaitGroup
	wg.Add(len(signerInfo.cosigner) - 1)
	go signerInfo.OtherNodeStartSign(&wg, &StartSignInfo{
		WorkId:   signerInfo.workId,
		Cosigner: signerInfo.cosigner,
		HashMsg:  signerInfo.hashMsg,
	})
	wg.Wait()
	return signerInfo.DoSignRound1()
}

func (s *SignOperator) TriggeredToStartSign(workId string, cosigner []uint32, hash []byte) error {
	s.cond.RLock()
	defer s.cond.RUnlock()
	var err error
	signerInfo := &SignerInfo{
		id:               s.id,
		workId:           workId,
		cosigner:         cosigner,
		hashMsg:          hash,
		participantAddrs: s.participantAddrs,
		round1RecvChan:   make(chan SignRound1Recv, s.threshold+1),
		round2RecvChan:   make(chan SignRound2Recv, s.threshold+1),
		round3RecvChan:   make(chan SignRound3Recv, s.threshold+1),
		round4RecvChan:   make(chan SignRound4Recv, s.threshold+1),
		round5RecvChan:   make(chan SignRound5Recv, s.threshold+1),
		round6RecvChan:   make(chan SignRound6Recv, s.threshold+1),
	}
	signerInfo.signer, err = GetDkgOperator().NewSigner(cosigner)
	if err != nil {
		return err
	}
	s.signerMap[workId] = signerInfo
	go func() {
		_, err := signerInfo.DoSignRound1()
		if err != nil {
			glog.Error(err)
		}
	}()
	return nil
}

func (si *SignerInfo) DoSignRound1() (*curves.EcdsaSignature, error) {
	round1, p2p, err := si.signer.SignRound1()
	if err != nil {
		return nil, err
	}
	signerOut := make(map[uint32]*ptcpt.Round1Bcast, len(si.cosigner)-1)
	r1P2pIn := make(map[uint32]*ptcpt.Round1P2PSend, len(si.cosigner)-1)

	go si.SendToOtherRound1Out(round1, p2p)

	for {
		if len(signerOut) == (len(si.cosigner) - 1) {
			break
		}
		select {
		case recv := <-si.round1RecvChan:
			if recv.WorkId != si.workId {
				continue
			}
			signerOut[recv.Id] = recv.SignerOut
			r1P2pIn[recv.Id] = recv.R1P2PSend
		case <-time.After(SignRoundWait * time.Second):
			glog.Error("等待SignerRound1Recv通道阻塞超时")
			return nil, errors.New("Sign Round1 Receive Wait timeout")
		}
	}
	return si.DoSignRound2(signerOut, r1P2pIn)
}

func (si *SignerInfo) DoSignRound2(params map[uint32]*ptcpt.Round1Bcast, p2p map[uint32]*ptcpt.Round1P2PSend) (*curves.EcdsaSignature, error) {
	round2Out, err := si.signer.SignRound2(params, p2p)
	if err != nil {
		return nil, err
	}
	r3In := make(map[uint32]*ptcpt.P2PSend, len(si.cosigner)-1)
	go si.SendToOtherRound2Out(round2Out)

	for {
		if len(r3In) == (len(si.cosigner) - 1) {
			break
		}
		select {
		case recv := <-si.round2RecvChan:
			if recv.WorkId != si.workId {
				continue
			}
			r3In[recv.Id] = recv.P2pSend
		case <-time.After(SignRoundWait * time.Second):
			glog.Error("等待SignerRound2Recv通道阻塞超时")
			return nil, errors.New("Sign Round2 Receive Wait timeout")
		}
	}
	return si.DoSignRound3(r3In)
}

func (si *SignerInfo) DoSignRound3(r3In map[uint32]*ptcpt.P2PSend) (*curves.EcdsaSignature, error) {
	round3Out, err := si.signer.SignRound3(r3In)
	if err != nil {
		return nil, err
	}
	r4In := make(map[uint32]*ptcpt.Round3Bcast, len(si.cosigner)-1)
	go si.SendToOtherRound3Out(round3Out)

	for {
		if len(r4In) == (len(si.cosigner) - 1) {
			break
		}
		select {
		case recv := <-si.round3RecvChan:
			if recv.WorkId != si.workId {
				continue
			}
			r4In[recv.Id] = recv.Round3Bcast
		case <-time.After(SignRoundWait * time.Second):
			glog.Error("等待SignerRound3Recv通道阻塞超时")
			return nil, errors.New("Sign Round3 Receive Wait timeout")
		}
	}

	return si.DoSignRound4(r4In)
}

func (si *SignerInfo) DoSignRound4(r4In map[uint32]*ptcpt.Round3Bcast) (*curves.EcdsaSignature, error) {
	round4Out, err := si.signer.SignRound4(r4In)
	if err != nil {
		return nil, err
	}
	r5In := make(map[uint32]*ptcpt.Round4Bcast, len(si.cosigner)-1)
	go si.SendToOtherRound4Out(round4Out)

	for {
		if len(r5In) == (len(si.cosigner) - 1) {
			break
		}
		select {
		case recv := <-si.round4RecvChan:
			if recv.WorkId != si.workId {
				continue
			}
			r5In[recv.Id] = recv.Round4Bcast
		case <-time.After(SignRoundWait * time.Second):
			glog.Error("等待SignerRound4Recv通道阻塞超时")
			return nil, errors.New("Sign Round4 Receive Wait timeout")
		}
	}
	return si.DoSignRound5(r5In)
}

func (si *SignerInfo) DoSignRound5(r5In map[uint32]*ptcpt.Round4Bcast) (*curves.EcdsaSignature, error) {
	round5Out, round5P2p, err := si.signer.SignRound5(r5In)
	if err != nil {
		return nil, err
	}
	r6In := make(map[uint32]*ptcpt.Round5Bcast, len(si.cosigner)-1)
	r6P2p := make(map[uint32]*ptcpt.Round5P2PSend, len(si.cosigner)-1)
	go si.SendToOtherRound5Out(round5Out, round5P2p)

	for {
		if len(r6In) == (len(si.cosigner) - 1) {
			break
		}
		select {
		case recv := <-si.round5RecvChan:
			if recv.WorkId != si.workId {
				continue
			}
			r6In[recv.Id] = recv.Round5Bcast
			r6P2p[recv.Id] = recv.Round5P2pSend
		case <-time.After(SignRoundWait * time.Second):
			glog.Error("等待SignerRound5Recv通道阻塞超时")
			return nil, errors.New("Sign Round5 Receive Wait timeout")
		}
	}
	return si.DoSignRound6(r6In, r6P2p)
}

func (si *SignerInfo) DoSignRound6(r6In map[uint32]*ptcpt.Round5Bcast, r6P2p map[uint32]*ptcpt.Round5P2PSend) (*curves.EcdsaSignature, error) {
	r6Bcast, err := si.signer.SignRound6Full(si.hashMsg, r6In, r6P2p)
	if err != nil {
		return nil, err
	}
	signOutIn := make(map[uint32]*ptcpt.Round6FullBcast, len(si.cosigner)-1)
	go si.SendToOtherRound6Out(r6Bcast)

	for {
		if len(signOutIn) == (len(si.cosigner) - 1) {
			break
		}
		select {
		case recv := <-si.round6RecvChan:
			if recv.WorkId != si.workId {
				continue
			}
			signOutIn[recv.Id] = recv.Round6FullBcast
		case <-time.After(SignRoundWait * time.Second):
			glog.Error("等待SignerRound6Recv通道阻塞超时")
			return nil, errors.New("Sign Round6 Receive Wait timeout")
		}
	}
	return si.DoGetSignOutput(signOutIn)
}

func (si *SignerInfo) DoGetSignOutput(signOutIn map[uint32]*ptcpt.Round6FullBcast) (*curves.EcdsaSignature, error) {
	defer func() {
		// 本次签名完成，及时将相关数据清除
		delete(GetSignOperator().signerMap, si.workId)
	}()
	return si.signer.SignOutput(signOutIn)
}

func (si *SignerInfo) OtherNodeStartSign(wg *sync.WaitGroup, info *StartSignInfo) {
	signOp := GetSignOperator()
	for _, nodeId := range si.cosigner {
		if signOp.id == nodeId {
			continue
		}
		nodeAddr := signOp.participantAddrs[nodeId-1]
		url := GetOtherStartSignUrl(nodeAddr)
		nodeId := nodeId
		go func() {
			err := DoPostWithoutRespData(url, *info)
			wg.Done()
			if err != nil {
				glog.Errorf("%d向%d触发开始签名失败: %v\n", signOp.id, nodeId, err.Error())
				return
			}
		}()
	}
}

func (s *SignOperator) UpdateInfo(id, threshold, total uint32, participantAddrs []string, otherParticipants []uint32) {
	s.cond.Lock()
	defer s.cond.Unlock()
	s.id = id
	s.threshold = threshold
	s.total = total
	s.participantAddrs = participantAddrs
	s.otherParticipants = otherParticipants
}
