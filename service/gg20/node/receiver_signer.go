package node

import (
	ptcpt "github.com/coinbase/kryptology/pkg/tecdsa/gg20/participant"
	"github.com/pkg/errors"
)

type SignRound1Recv struct {
	Id        uint32               `json:"id"`
	WorkId    string               `json:"workId"`
	SignerOut *ptcpt.Round1Bcast   `json:"signerOut"`
	R1P2PSend *ptcpt.Round1P2PSend `json:"r1P2PSend"`
}

type SignRound2Recv struct {
	Id      uint32         `json:"id"`
	WorkId  string         `json:"workId"`
	P2pSend *ptcpt.P2PSend `json:"p2PSend"`
}

type SignRound3Recv struct {
	Id          uint32             `json:"id"`
	WorkId      string             `json:"workId"`
	Round3Bcast *ptcpt.Round3Bcast `json:"round3Bcast"`
}

type SignRound4Recv struct {
	Id          uint32             `json:"id"`
	WorkId      string             `json:"workId"`
	Round4Bcast *ptcpt.Round4Bcast `json:"round4Bcast"`
}
type SignRound5Recv struct {
	Id            uint32               `json:"id"`
	WorkId        string               `json:"workId"`
	Round5Bcast   *ptcpt.Round5Bcast   `json:"round5Bcast"`
	Round5P2pSend *ptcpt.Round5P2PSend `json:"round5P2PSend"`
}

type SignRound6Recv struct {
	Id              uint32                 `json:"id"`
	WorkId          string                 `json:"workId"`
	Round6FullBcast *ptcpt.Round6FullBcast `json:"round6FullBcast"`
}

// RecvAskCandidateInfo 接收到需要签名的请求，将workId与自身id一起返回
func (s *SignOperator) RecvAskCandidateInfo(workId string) *CandidateInfo {
	return &CandidateInfo{Id: s.id, WorkId: workId}
}

// DoSignRound1Recv 接收其它节点广播发来的SignRound1的结果，将其放入对应WorkId下round1的数据接收通道中
func (s *SignOperator) DoSignRound1Recv(info *SignRound1Recv) error {
	signerInfo, exist := s.signerMap[info.WorkId]
	if !exist {
		return errors.Errorf("当前WorkId %s 的signer不存在", info.WorkId)
	}
	signerInfo.round1RecvChan <- *info
	return nil
}

// DoSignRound2Recv 类似DoSignRound1Recv，接收数据并推入对应通道，下同
func (s *SignOperator) DoSignRound2Recv(info *SignRound2Recv) error {
	signerInfo, exist := s.signerMap[info.WorkId]
	if !exist {
		return errors.Errorf("当前WorkId %s 的signer不存在", info.WorkId)
	}
	signerInfo.round2RecvChan <- *info
	return nil
}

func (s *SignOperator) DoSignRound3Recv(info *SignRound3Recv) error {
	signerInfo, exist := s.signerMap[info.WorkId]
	if !exist {
		return errors.Errorf("当前WorkId %s 的signer不存在", info.WorkId)
	}
	signerInfo.round3RecvChan <- *info
	return nil
}

func (s *SignOperator) DoSignRound4Recv(info *SignRound4Recv) error {
	signerInfo, exist := s.signerMap[info.WorkId]
	if !exist {
		return errors.Errorf("当前WorkId %s 的signer不存在", info.WorkId)
	}
	signerInfo.round4RecvChan <- *info
	return nil
}

func (s *SignOperator) DoSignRound5Recv(info *SignRound5Recv) error {
	signerInfo, exist := s.signerMap[info.WorkId]
	if !exist {
		return errors.Errorf("当前WorkId %s 的signer不存在", info.WorkId)
	}
	signerInfo.round5RecvChan <- *info
	return nil
}

func (s *SignOperator) DoSignRound6Recv(info *SignRound6Recv) error {
	signerInfo, exist := s.signerMap[info.WorkId]
	if !exist {
		return errors.Errorf("当前WorkId %s 的signer不存在", info.WorkId)
	}
	signerInfo.round6RecvChan <- *info
	return nil
}
