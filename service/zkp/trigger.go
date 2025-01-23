package zkp

import (
	"github.com/coinbase/kryptology/service/gg20/node"
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

	operator := node.GetDkgOperator()
	operator.UpdateClusterInfo(GetZkManager().nodeAddress, children) // 更新集群信息，并设置其状态为不可用

	zkLock := NewZkLock(GetZkManager().GetConn(), DkgLockPath)

	isLeader, err := zkLock.AcquireNotWatch()

	defer func() {
		if err := zkLock.Release(); err != nil {
			glog.Error("释放锁失败: ", err)
		}
	}()

	if err != nil {
		glog.Info("获取分布式锁失败: ", err)
		return err
	} else if !isLeader {
		glog.Info("未获取到分布式锁：", zkLock.ownNodePath)
	} else if isLeader {
		executeDKG(GetZkManager().nodeAddress, children)
	}

	return nil
}

func executeDKG(myaddress string, nodes []string) {
	glog.Info("Executing DKG...")
	err := node.GetDkgOperator().StartDkg()
	if err != nil {
		glog.Errorf("执行Dkg失败：%v", err.Error())
		return
	}
	glog.Info("DKG completed.")
}
