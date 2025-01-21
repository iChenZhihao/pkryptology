package node

import (
	"fmt"
	"testing"
)

func TestUpdateClusterInfo(t *testing.T) {
	addrs := []string{"192.168.0.105:8081", "192.168.0.105:8082", "192.168.0.105:8080"}
	GetDkgOperator().UpdateClusterInfo("192.168.0.105:8081", addrs)
	operator := GetDkgOperator()
	fmt.Printf("generator:%v   otherParticipants:%v\n", generator, operator.otherParticipants)
}
