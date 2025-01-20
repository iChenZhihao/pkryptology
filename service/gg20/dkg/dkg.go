package dkg

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"github.com/btcsuite/btcd/btcec"
	"github.com/coinbase/kryptology/pkg/core"
	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/paillier"
	v1 "github.com/coinbase/kryptology/pkg/sharing/v1"
	part "github.com/coinbase/kryptology/pkg/tecdsa/gg20/participant"
	//"github.com/coinbase/kryptology/service/global"
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

	ecdsaVerifier = func(verKey *curves.EcPoint, hash []byte, sig *curves.EcdsaSignature) bool {
		pk := &ecdsa.PublicKey{
			Curve: verKey.Curve,
			X:     verKey.X,
			Y:     verKey.Y,
		}
		return ecdsa.Verify(pk, hash, sig.R, sig.S)
	}
	curve = elliptic.P256()
)

type DkgOperator struct {
	nodeAddress       string
	participantAddrs  []string
	participant       *part.DkgParticipant
	dkgResult         *part.DkgResult
	id                uint32
	threshold         uint32
	total             uint32
	ChanRecvRound1    chan DkgRound1Recv
	ChanRecvRound2    chan DkgRound2Recv
	ChanRecvRound3    chan DkgRound3Recv
	available         bool
	otherParticipants []uint32
	cond              sync.RWMutex
}

type DkgRound1Recv struct {
	Id uint32 `json:"id"`
	//Bcastj kdkg.Round1Bcast   `json:"bcastj"`
	//P2pj   kdkg.Round1P2PSend `json:"p2pj"`
	Round1Bcast *part.DkgRound1Bcast `json:"round1Bcast"`
}

type DkgRound2Recv struct {
	Id uint32 `json:"id"`
	//Round2Outj kdkg.Round2Bcast `json:"round2Outj"`
	Decommitment *core.Witness   `json:"decommitment"`
	ShamirShare  *v1.ShamirShare `json:"shamirShare"`
}

type DkgRound3Recv struct {
	Id       uint32            `json:"id"`
	PsfProof paillier.PsfProof `json:"psfProof"`
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
	d.total = uint32(len(participants))
	d.threshold = (d.total / 2) + 1
	if d.ChanRecvRound1 != nil {
		close(d.ChanRecvRound1)
	}
	d.ChanRecvRound1 = make(chan DkgRound1Recv, len(participants)*2) //预留一些缓冲
	if d.ChanRecvRound2 != nil {
		close(d.ChanRecvRound2)
	}
	d.ChanRecvRound2 = make(chan DkgRound2Recv, len(participants)*2) //预留一些缓冲
	if d.ChanRecvRound3 != nil {
		close(d.ChanRecvRound3)
	}
	d.ChanRecvRound3 = make(chan DkgRound3Recv, len(participants)*2)

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
	glog.Infof("MyId:%v, otherParticipants: %v\n", d.id, d.otherParticipants)
	glog.Info("update dkgOperator info success!!!~~~~~~~~~")
}

func (d *DkgOperator) StartDkg() error {
	d.cond.RLock()
	defer d.cond.RUnlock()
	var wg sync.WaitGroup
	wg.Add(len(d.participantAddrs) - 1)
	go d.OtherNodeStartDkg(&wg)

	err := d.DoDkgRound1()
	wg.Wait()
	if err != nil {
		glog.Errorf("")
	}
	return nil
}

func (d *DkgOperator) TriggeredToStartRound1() {
	d.cond.RLock()
	defer d.cond.RUnlock()
	go func() {
		err := d.DoDkgRound1()
		if err != nil {
			glog.Error(nil)
		}
	}()
}

// DoDkgRound1 初始化participant后执行Round1，完成并广播、接收数据后，执行Round2
func (d *DkgOperator) DoDkgRound1() error {
	glog.Infof("DoDkgRound1 otherParticipants: %v\n", d.otherParticipants)
	d.participant = part.NewDkgParticipant(curve, d.id, d.threshold, d.total)

	//nodeBcast, nodeP2psend, _ := d.participant.Round1(secret)
	dkgRound1, err := d.participant.DkgRound1(d.threshold, d.total)
	if err != nil {
		glog.Errorf("")
	}

	dkgR1Outs := make(map[uint32]*part.DkgRound1Bcast, d.total)
	//bcast := make(map[uint32]kdkg.Round1Bcast)
	//nodeP2p := make(map[uint32]*kdkg.Round1P2PSendPacket)
	// 当前节点i 的 nodeP2p[j] 为其它节点（例如j）的发来的nodeP2pJ[i]

	// 给其它节点广播nodeP2psend
	toSend := &DkgRound1Recv{Id: d.id, Round1Bcast: dkgRound1}
	go d.SendToOtherNodesDkgRound1Out(*toSend)

	dkgR1Outs[d.id] = dkgRound1

	// 等待其它节点发来的数据
	for {
		if len(dkgR1Outs) == len(d.participantAddrs) {
			break
		}
		select {
		case recv := <-d.ChanRecvRound1:
			//bcast[recv.Id] = recv.Bcastj
			//nodeP2p[recv.Id] = recv.P2pj[d.id]
			dkgR1Outs[recv.Id] = recv.Round1Bcast
		case <-time.After(200 * time.Second):
			glog.Error("等待Round1Recv通道阻塞超时")
			return errors.New("Dkg Round1 Receive Wait timeout")
		}
	}
	err = d.DoDkgRound2(dkgR1Outs)

	return nil

}

// DoDkgRound2 执行Round2，完成并广播、接收数据后，执行Round3与4
func (d *DkgOperator) DoDkgRound2(dkgR1Outs map[uint32]*part.DkgRound1Bcast) error {
	//nodeRound2Out, err := d.participant.Round2(bcast, p2p)
	dkgR2Bcast, _, err := d.participant.DkgRound2(dkgR1Outs)
	if err != nil {
		glog.Error("Participant执行Round2出错：", err)
		return err
	}
	//round3Input := make(map[uint32]kdkg.Round2Bcast)
	decommitments := make(map[uint32]*core.Witness, d.total)
	decommitments[d.id] = dkgR2Bcast.Di
	shamirMap := make(map[uint32]*v1.ShamirShare, d.total)
	X := d.participant.GetShamirShamirX()
	shamirMap[d.id] = X[d.id-1]

	toSend := &DkgRound2Recv{Id: d.id, Decommitment: dkgR2Bcast.Di}
	go d.SendToOtherNodesDkgRound2Out(*toSend, X)

	for {
		if len(decommitments) == len(d.participantAddrs) {
			break
		}
		select {
		case recv := <-d.ChanRecvRound2:
			decommitments[recv.Id] = recv.Decommitment
			shamirMap[recv.Id] = recv.ShamirShare
		case <-time.After(20 * time.Second):
			glog.Error("等待Round2Recv通道阻塞超时")
			return errors.New("Dkg Round2 Receive Wait timeout")
		}
	}

	err = d.DoDkgRound3(decommitments, shamirMap)
	if err != nil {
		return err
	}
	return nil
}

func (d *DkgOperator) DoDkgRound3(decommitments map[uint32]*core.Witness, shamirMap map[uint32]*v1.ShamirShare) error {
	var err error
	dkgR3OutMap := make(map[uint32]paillier.PsfProof, d.total)
	dkgR3OutMap[d.id], err = d.participant.DkgRound3(decommitments, shamirMap)
	if err != nil {
		glog.Error("Participant执行Round3出错：", err)
		return err
	}
	toSend := &DkgRound3Recv{Id: d.id, PsfProof: dkgR3OutMap[d.id]}
	go d.SendToOtherNodesDkgRound3Out(*toSend)

	for {
		if len(dkgR3OutMap) == len(d.participantAddrs) {
			break
		}
		select {
		case recv := <-d.ChanRecvRound3:
			dkgR3OutMap[recv.Id] = recv.PsfProof
		case <-time.After(30 * time.Second):
			glog.Error("等待Round3Recv通道阻塞超时")
			return errors.New("Dkg Round3 Receive Wait timeout")
		}
	}

	err = d.DoDkgRound4(dkgR3OutMap)
	if err != nil {
		return err
	}
	return nil
}

func (d *DkgOperator) DoDkgRound4(psfProof map[uint32]paillier.PsfProof) error {
	round4, err := d.participant.DkgRound4(psfProof)
	if err != nil {
		glog.Error("Participant执行Round4出错：", err)
		return err
	}
	d.dkgResult = round4
	d.available = true
	return nil
}

func (d *DkgOperator) OtherNodeStartDkg(group *sync.WaitGroup) {
	for _, nodeId := range d.otherParticipants {
		url := OtherNodeStartDkgUrl(d.participantAddrs[nodeId-1])
		nodeId := nodeId
		go func() {
			err := DoSendStartDkg(url)
			group.Done()
			if err != nil {
				glog.Errorf("%d向%d触发启动DKG失败: %v\n", d.id, nodeId, err.Error())
				return
			}
		}()
	}
}
