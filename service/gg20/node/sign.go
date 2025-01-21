package node

import (
	"fmt"
	"github.com/coinbase/kryptology/pkg/core/curves"
	ptcpt "github.com/coinbase/kryptology/pkg/tecdsa/gg20/participant"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"sync"
	"time"
)

var (
	signOnce     sync.Once
	signOperator *SignOperator
)

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
	workId            string
	cosigner          []uint32
	candidateInfoChan chan *CandidateInfo
	signer            *ptcpt.Signer
}

func (s *SignOperator) SignMsg(plainMsg []byte) (*curves.EcdsaSignature, error) {
	s.cond.RLock()
	s.cond.RUnlock()
	if !GetDkgOperator().IsAvailable() {
		return nil, fmt.Errorf("签名节点集群暂不可用")
	}
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	signerInfo := &SignerInfo{
		workId:            newUUID.String(),
		cosigner:          []uint32{s.id},
		candidateInfoChan: make(chan *CandidateInfo, s.total+1),
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
		case <-time.After(15 * time.Second):
			return nil, fmt.Errorf("%d节点等待确定其它cosigner超时", s.id)
		}
	}
	signerInfo.signer, err = GetDkgOperator().NewSigner(signerInfo.cosigner)
	if err != nil {
		return nil, err
	}

	glog.Info("确定了cosigner：", s.signerMap[signerInfo.workId].cosigner)

	return nil, fmt.Errorf("暂时没开发好")
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
