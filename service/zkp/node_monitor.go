package zkp

import (
	"github.com/coinbase/kryptology/service/global"
	"github.com/go-zookeeper/zk"
	"github.com/golang/glog"
	"sync"
	"time"
)

var (
	stableTimer *time.Timer // 稳定窗口计时器
	mutex       sync.Mutex  // 互斥锁，确保计时器的安全操作
)

// MonitorNodeChanges 监听节点数量变化
func (m *ZkManager) MonitorNodeChanges() {
	for {
		children, _, ch, err := m.GetConn().ChildrenW(ServerRegisterPath)
		if err != nil {
			glog.Errorf("获取子节点失败: %v", err)
		}

		nodeCount := len(children)
		glog.Infof("当前节点数: %d\n", nodeCount)

		go m.updateStabilityTimer(nodeCount)

		select {
		case event := <-ch:
			if event.Type == zk.EventNodeChildrenChanged {
				continue
			}
		}
	}
}

// 更新稳定窗口计时器
func (m *ZkManager) updateStabilityTimer(nodeCount int) {
	mutex.Lock()
	defer mutex.Unlock()

	// 如果节点数发生变化，重置计时器
	if nodeCount != m.lastNodeCount {
		m.lastNodeCount = nodeCount

		// 重置计时器
		if stableTimer != nil {
			stableTimer.Stop()
		}

		// 创建新的计时器
		stableTimer = time.AfterFunc(global.Config.GG20.GetStableTimeWindow(), func() {
			// 计时器到期时触发DKG
			TriggerDKG(nodeCount)
		})
		glog.Infof("节点数变动，重置稳定窗口计时器为 %v 毫秒 \n", global.Config.GG20.GetStableTimeWindow().Milliseconds())
	}
}
