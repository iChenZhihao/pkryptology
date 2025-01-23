package node

import (
	"github.com/coinbase/kryptology/pkg/core"
	v1 "github.com/coinbase/kryptology/pkg/sharing/v1"
	"github.com/golang/glog"
)

// SendToOtherNodesDkgRound1Out 向其它节点广播本节点Round1的结果
func (d *DkgOperator) SendToOtherNodesDkgRound1Out(toSend DkgRound1Recv) {
	for _, nodeId := range d.otherParticipants {
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
}

// SendToOtherNodesDkgRound2Out 向其它节点广播本节点Round2的结果
func (d *DkgOperator) SendToOtherNodesDkgRound2Out(witness *core.Witness, X []*v1.ShamirShare) {
	for _, nodeId := range d.otherParticipants {
		url := GetDkgRound2BcastUrl(d.participantAddrs[nodeId-1])
		nodeId := nodeId
		go func() {
			// 并发地发送请求时，使用同一个对象可能导致协程间相互影响其中的数据，因此需要在协程内部创建数据对象并发送
			toSend := &DkgRound2Recv{Id: d.id, Decommitment: witness, ShamirShare: X[nodeId-1]}
			err := DoSendBroadcastRound2(url, *toSend)
			if err != nil {
				glog.Errorf("%d向%d广播Round2结果失败: %v\n", d.id, nodeId, err)
				return
			}
		}()
	}
}

func (d *DkgOperator) SendToOtherNodesDkgRound3Out(toSend DkgRound3Recv) {
	for _, nodeId := range d.otherParticipants {
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
