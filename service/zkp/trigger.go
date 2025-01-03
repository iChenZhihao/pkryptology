package zkp

import (
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

var (
	nodeCountDontMatch = errors.New("Node Count Don't Match")
)

func TriggerDKG(nodeCount int) error {
	glog.Infof("开启DKG流程, 节点数：%d个\n", nodeCount)

	children, _, err := GetZkManager().GetConn().Children(ServerRegisterPath)
	if err != nil {
		glog.Errorf("获取子节点失败: %v", err)
		return nil
	}
	if len(children) != nodeCount {
		glog.Errorf("")
		return nodeCountDontMatch
	}

	for index, nodeAddress := range children {
		glog.Infof("节点:%d号，地址:%s\n", index, nodeAddress)
	}
	return nil
}
