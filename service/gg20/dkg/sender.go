package dkg

import (
	v1 "github.com/coinbase/kryptology/pkg/sharing/v1"
	"github.com/golang/glog"
)

// SendToOtherNodesDkgRound1Out 向其它节点广播本节点Round1的结果
func (d *DkgOperator) SendToOtherNodesDkgRound1Out(toSend DkgRound1Recv) {
	//var wg sync.WaitGroup
	for _, nodeId := range d.otherParticipants {
		//wg.Add(1)
		//if d.nodeAddress == addr {
		//	continue
		//}
		url := GetDkgRound1BcastUrl(d.participantAddrs[nodeId-1])
		nodeId := nodeId
		go func() {
			err := DoSendBroadcastRound1(url, toSend)
			if err != nil {
				glog.Errorf("%d向%d广播Round1结果失败: %v\n", d.id, nodeId, err.Error())
				return
			}
		}()
	}
	//wg.Wait()
}

// SendToOtherNodesDkgRound2Out 向其它节点广播本节点Round2的结果
func (d *DkgOperator) SendToOtherNodesDkgRound2Out(toSend DkgRound2Recv, X []*v1.ShamirShare) {
	//var wg sync.WaitGroup
	for _, nodeId := range d.otherParticipants {
		//wg.Add(1)
		url := GetDkgRound2BcastUrl(d.participantAddrs[nodeId-1])
		toSend.ShamirShare = X[nodeId-1]
		nodeId := nodeId
		go func() {
			err := DoSendBroadcastRound2(url, toSend)
			if err != nil {
				glog.Errorf("%d向%d广播Round2结果失败: %v\n", d.id, nodeId, err)
				return
			}
		}()
	}
	//wg.Wait()
}

func (d *DkgOperator) SendToOtherNodesDkgRound3Out(toSend DkgRound3Recv) {
	//var wg sync.WaitGroup
	for _, nodeId := range d.otherParticipants {
		//wg.Add(1)
		url := GetDkgRound3BcastUrl(d.participantAddrs[nodeId-1])
		nodeId := nodeId
		go func() {
			err := DoSendBroadcastRound3(url, toSend)
			if err != nil {
				glog.Errorf("%d向%d广播Round3结果失败: %v\n", d.id, nodeId, err)
				return
			}
		}()
	}
}
