package node

import (
	ptcpt "github.com/coinbase/kryptology/pkg/tecdsa/gg20/participant"
	"github.com/golang/glog"
)

type StartSignInfo struct {
	WorkId   string   `json:"workId"`
	Cosigner []uint32 `json:"cosigner"`
	HashMsg  []byte   `json:"hashMsg"`
}

func (si *SignerInfo) SendToOtherRound1Out(round1Bcast *ptcpt.Round1Bcast, p2pSend map[uint32]*ptcpt.Round1P2PSend) {
	for _, nodeId := range si.cosigner {
		if nodeId == si.id {
			continue
		}
		nodeAddr := si.participantAddrs[nodeId-1]
		url := GetSignRound1BcastUrl(nodeAddr)
		nodeId := nodeId
		go func() {
			data := &SignRound1Recv{
				Id:        si.id,
				WorkId:    si.workId,
				R1P2PSend: p2pSend[nodeId],
				SignerOut: round1Bcast,
			}
			err := DoPostWithoutRespData(url, data)
			if err != nil {
				glog.Errorf("%d向%d广播SignRound1结果失败: %v\n", si.id, nodeId, err.Error())
				return
			}
		}()
	}
}

func (si *SignerInfo) SendToOtherRound2Out(p2pSend map[uint32]*ptcpt.P2PSend) {
	for _, nodeId := range si.cosigner {
		if nodeId == si.id {
			continue
		}
		nodeAddr := si.participantAddrs[nodeId-1]
		url := GetSignRound2BcastUrl(nodeAddr)
		nodeId := nodeId
		go func() {
			data := &SignRound2Recv{
				Id:      si.id,
				WorkId:  si.workId,
				P2pSend: p2pSend[nodeId],
			}
			err := DoPostWithoutRespData(url, data)
			if err != nil {
				glog.Errorf("%d向%d广播SignRound2结果失败: %v\n", si.id, nodeId, err.Error())
				return
			}
		}()
	}
}

func (si *SignerInfo) SendToOtherRound3Out(round3Bcast *ptcpt.Round3Bcast) {
	for _, nodeId := range si.cosigner {
		if nodeId == si.id {
			continue
		}
		nodeAddr := si.participantAddrs[nodeId-1]
		url := GetSignRound3BcastUrl(nodeAddr)
		nodeId := nodeId
		go func() {
			data := &SignRound3Recv{
				Id:          si.id,
				WorkId:      si.workId,
				Round3Bcast: round3Bcast,
			}
			err := DoPostWithoutRespData(url, data)
			if err != nil {
				glog.Errorf("%d向%d广播SignRound3结果失败: %v\n", si.id, nodeId, err.Error())
				return
			}
		}()
	}
}

func (si *SignerInfo) SendToOtherRound4Out(round4Bcast *ptcpt.Round4Bcast) {
	for _, nodeId := range si.cosigner {
		if nodeId == si.id {
			continue
		}
		nodeAddr := si.participantAddrs[nodeId-1]
		url := GetSignRound4BcastUrl(nodeAddr)
		nodeId := nodeId
		go func() {
			data := &SignRound4Recv{
				Id:          si.id,
				WorkId:      si.workId,
				Round4Bcast: round4Bcast,
			}
			err := DoPostWithoutRespData(url, data)
			if err != nil {
				glog.Errorf("%d向%d广播SignRound4结果失败: %v\n", si.id, nodeId, err.Error())
				return
			}
		}()
	}
}

func (si *SignerInfo) SendToOtherRound5Out(round5Bcast *ptcpt.Round5Bcast, r5P2p map[uint32]*ptcpt.Round5P2PSend) {
	for _, nodeId := range si.cosigner {
		if nodeId == si.id {
			continue
		}
		nodeAddr := si.participantAddrs[nodeId-1]
		url := GetSignRound5BcastUrl(nodeAddr)
		nodeId := nodeId
		go func() {
			data := &SignRound5Recv{
				Id:            si.id,
				WorkId:        si.workId,
				Round5Bcast:   round5Bcast,
				Round5P2pSend: r5P2p[nodeId],
			}
			err := DoPostWithoutRespData(url, data)
			if err != nil {
				glog.Errorf("%d向%d广播SignRound5结果失败: %v\n", si.id, nodeId, err.Error())
				return
			}
		}()
	}
}

func (si *SignerInfo) SendToOtherRound6Out(round6Bcast *ptcpt.Round6FullBcast) {
	for _, nodeId := range si.cosigner {
		if nodeId == si.id {
			continue
		}
		nodeAddr := si.participantAddrs[nodeId-1]
		url := GetSignRound6BcastUrl(nodeAddr)
		nodeId := nodeId
		go func() {
			data := &SignRound6Recv{
				Id:              si.id,
				WorkId:          si.workId,
				Round6FullBcast: round6Bcast,
			}
			err := DoPostWithoutRespData(url, data)
			if err != nil {
				glog.Errorf("%d向%d广播SignRound4结果失败: %v\n", si.id, nodeId, err.Error())
				return
			}
		}()
	}
}
