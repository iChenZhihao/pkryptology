package dkg

import (
	"encoding/json"
	"fmt"
	"github.com/coinbase/kryptology/pkg/core/curves"
	kdkg "github.com/coinbase/kryptology/pkg/dkg/gennaro"
	"testing"
)

func TestDkgRound1RecvSerialize(t *testing.T) {
	participant, err := kdkg.NewParticipant(1,
		2,
		generator,
		curves.NewK256Scalar(),
		2, 3)
	if err != nil {
		t.Failed()
	}
	round1, send, err := participant.Round1(nil)
	toSend := &DkgRound1Recv{Id: 1, Bcastj: round1, P2pj: send}
	jsonData, err := json.Marshal(toSend)
	if err != nil {
		t.Failed()
	}
	Deserialize := &DkgRound1Recv{}
	err = json.Unmarshal(jsonData, Deserialize)
	fmt.Printf("%v\n", Deserialize)
}

func TestUpdateClusterInfo(t *testing.T) {
	addrs := []string{"192.168.0.105:8081", "192.168.0.105:8082", "192.168.0.105:8080"}
	GetDkgOperator().UpdateClusterInfo("192.168.0.105:8081", addrs)
	operator := GetDkgOperator()
	fmt.Printf("generator:%v   otherParticipants:%v\n", generator, operator.otherParticipants)
}
