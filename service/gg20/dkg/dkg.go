package dkg

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/coinbase/kryptology/pkg/core/curves"
	kdkg "github.com/coinbase/kryptology/pkg/dkg/gennaro"
	"github.com/coinbase/kryptology/service/global"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"math/big"
	"sync"
	"time"
)

const DkgSecretKey = "secret"

var (
	mu           sync.Mutex
	Operator     *DkgOperator
	once         sync.Once
	generator, _ = curves.NewScalarBaseMult(btcec.S256(), big.NewInt(4765))
)

type DkgOperator struct {
	nodeAddress       string
	participantAddrs  []string
	participant       *kdkg.Participant
	id                uint32
	threshold         uint32
	ChanRecvRound1    chan DkgRound1Recv
	ChanRecvRound2    chan DkgRound2Recv
	available         bool
	otherParticipants []uint32
	cond              sync.RWMutex
}

type DkgRound1Recv struct {
	Id     uint32             `json:"id"`
	Bcastj kdkg.Round1Bcast   `json:"bcastj"`
	P2pj   kdkg.Round1P2PSend `json:"p2pj"`
}

type DkgRound2Recv struct {
	Id         uint32           `json:"id"`
	Round2Outj kdkg.Round2Bcast `json:"round2Outj"`
}

func GetDkgOperator() *DkgOperator {
	once.Do(func() {
		Operator = &DkgOperator{}
	})
	return Operator
}

// UpdateClusterInfo 更新签名节点集群信息
func (d *DkgOperator) UpdateClusterInfo(nodeAddress string, participants []string) {
	//mu.Lock()
	//defer mu.Unlock()
	d.cond.Lock()
	defer d.cond.Unlock()
	d.available = false
	d.nodeAddress = nodeAddress
	d.participantAddrs = participants
	if d.ChanRecvRound1 != nil {
		close(d.ChanRecvRound1)
	}
	d.ChanRecvRound1 = make(chan DkgRound1Recv, len(participants)*2) //预留一些缓冲
	if d.ChanRecvRound2 != nil {
		close(d.ChanRecvRound2)
	}
	d.ChanRecvRound2 = make(chan DkgRound2Recv, len(participants)*2) //预留一些缓冲

	currentId := -1
	d.otherParticipants = make([]uint32, 0)
	for idx, addr := range d.participantAddrs {
		if d.nodeAddress == addr {
			currentId = idx
		} else {
			d.otherParticipants = append(d.otherParticipants, uint32(idx+1))
		}
	}
	if currentId < 0 {
		glog.Error("未能确定当前节点的id")
	}
	currentId += 1
	d.id = uint32(currentId)
	d.threshold = (uint32(len(d.participantAddrs)) / 2) + 1
	glog.Infof("MyId:%v, otherParticipants: %v\n", d.id, d.otherParticipants)
	glog.Info("update dkgOperator info success!!!~~~~~~~~~")
}

func (d *DkgOperator) StartDkg() error {
	d.cond.RLock()
	defer d.cond.RUnlock()
	var wg sync.WaitGroup
	wg.Add(len(d.participantAddrs) - 1)
	go d.OtherNodeStartDkg(&wg)

	err := d.DoDkgRound1(SecretStrToBytes(global.Config.GG20.Secret))
	wg.Wait()
	if err != nil {
		glog.Errorf("")
	}
	return nil
}

func (d *DkgOperator) TriggeredToStartRound1(secret []byte) {
	d.cond.RLock()
	defer d.cond.RUnlock()
	go func() {
		err := d.DoDkgRound1(secret)
		if err != nil {
			glog.Error(nil)
		}
	}()
}

// DoDkgRound1 初始化participant后执行Round1，完成并广播、接收数据后，执行Round2
func (d *DkgOperator) DoDkgRound1(secret []byte) error {
	d.cond.RLock()
	defer d.cond.RUnlock()
	glog.Infof("DoDkgRound1 otherParticipants: %v\n", d.otherParticipants)
	participant, err := kdkg.NewParticipant(d.id,
		d.threshold,
		generator,
		curves.NewK256Scalar(),
		d.otherParticipants...)
	if err != nil {
		glog.Errorf("创建Participant失败：%v\n", err.Error())
		return err
	}
	d.participant = participant

	nodeBcast, nodeP2psend, _ := d.participant.Round1(secret)

	bcast := make(map[uint32]kdkg.Round1Bcast)
	nodeP2p := make(map[uint32]*kdkg.Round1P2PSendPacket)
	// 当前节点i 的 nodeP2p[j] 为其它节点（例如j）的发来的nodeP2pJ[i]

	// 给其它节点广播nodeP2psend
	toSend := &DkgRound1Recv{Id: d.id, Bcastj: nodeBcast, P2pj: nodeP2psend}
	go d.SendToOtherNodesDkgRound1Out(*toSend)

	bcast[d.id] = nodeBcast

	// 等待其它节点发来的数据
	for {
		if len(bcast) == len(d.participantAddrs) {
			break
		}
		select {
		case recv := <-d.ChanRecvRound1:
			bcast[recv.Id] = recv.Bcastj
			nodeP2p[recv.Id] = recv.P2pj[d.id]
		case <-time.After(20 * time.Second):
			glog.Error("等待Round1Recv通道阻塞超时")
			return errors.New("Dkg Round1 Receive Wait timeout")
		}
	}
	err = d.DoDkgRound2(bcast, nodeP2p)

	return nil

}

// DoDkgRound2 执行Round2，完成并广播、接收数据后，执行Round3与4
func (d *DkgOperator) DoDkgRound2(bcast map[uint32]kdkg.Round1Bcast, p2p map[uint32]*kdkg.Round1P2PSendPacket) error {
	nodeRound2Out, err := d.participant.Round2(bcast, p2p)
	if err != nil {
		glog.Error("Participant执行Round2出错：", err)
		return err
	}
	round3Input := make(map[uint32]kdkg.Round2Bcast)
	round3Input[d.id] = nodeRound2Out

	toSend := &DkgRound2Recv{Id: d.id, Round2Outj: nodeRound2Out}
	go d.SendToOtherNodesDkgRound2Out(*toSend)

	for {
		if len(round3Input) == len(d.participantAddrs) {
			break
		}
		select {
		case recv := <-d.ChanRecvRound2:
			round3Input[recv.Id] = recv.Round2Outj
		case <-time.After(20 * time.Second):
			glog.Error("等待Round2Recv通道阻塞超时")
			return errors.New("Dkg Round2 Receive Wait timeout")
		}
	}

	err = d.doDkgRound3And4(round3Input)
	if err != nil {
		return err
	}
	return nil
}

// DoDkgRound3And4 执行Round3 与 Round4 （不依赖其它节点的结果，因此合并执行）
func (d *DkgOperator) doDkgRound3And4(bcast map[uint32]kdkg.Round2Bcast) error {
	round3Out, _, err := d.participant.Round3(bcast)
	if round3Out == nil {
		glog.Errorf("%d号节点执行dkg Round3结果为空\n", d.id)
	}
	if err != nil {
		glog.Errorf("%d号节点执行dkg Round3出错：%v \n", d.id, err)
		return err
	}
	round4Out, err := d.participant.Round4()
	if err != nil {
		glog.Error("%d号节点执行dkg Round4出错：%v \n", d.id, err)
		return err
	}
	glog.Infof("%d号节点的公钥：%v\n", d.id, round4Out)
	for keyId, value := range round4Out {
		glog.Infof("%d: %v\n", keyId, value)
	}
	d.available = true
	return nil
}

func (d *DkgOperator) OtherNodeStartDkg(group *sync.WaitGroup) {
	for _, nodeId := range d.otherParticipants {
		url := OtherNodeStartDkgUrl(d.participantAddrs[nodeId-1])
		nodeId := nodeId
		go func() {
			err := DoSendStartDkg(url, SecretStrToBase64Str(global.Config.GG20.Secret))
			group.Done()
			if err != nil {
				glog.Errorf("%d向%d触发启动DKG失败: %v\n", d.id, nodeId, err.Error())
				return
			}
		}()
	}
}

// RecvDkgRound1 接收其它节点（称之为j）发来的第一轮数据
func (d *DkgOperator) RecvDkgRound1(recv DkgRound1Recv) {
	d.ChanRecvRound1 <- recv
}

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

// RecvDkgRound2 接收其它节点（称之为j）发来的第一轮数据
func (d *DkgOperator) RecvDkgRound2(recv DkgRound2Recv) {
	d.ChanRecvRound2 <- recv
}

// SendToOtherNodesDkgRound2Out 向其它节点广播本节点Round2的结果
func (d *DkgOperator) SendToOtherNodesDkgRound2Out(toSend DkgRound2Recv) {
	//var wg sync.WaitGroup
	for _, nodeId := range d.otherParticipants {
		//wg.Add(1)
		url := GetDkgRound2BcastUrl(d.participantAddrs[nodeId-1])
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
