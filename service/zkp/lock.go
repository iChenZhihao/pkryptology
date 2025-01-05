package zkp

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// 用于在 TriggerDkg 时，基于Zk分布式锁，以防止DKG被多次执行

const DkgLockPath = "/gg20/dkg_lock"

type ZkLock struct {
	conn        *zk.Conn
	lockPath    string
	ownNodePath string
}

func NewZkLock(conn *zk.Conn, lockPath string) *ZkLock {
	return &ZkLock{
		conn:     conn,
		lockPath: lockPath,
	}
}

// Acquire 获取锁
func (z *ZkLock) Acquire() error {
	// 确保锁目录存在
	exists, _, err := z.conn.Exists(z.lockPath)
	if err != nil {
		return err
	}
	if !exists {
		_, err = z.conn.Create(z.lockPath, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}

	// 创建临时有序节点
	nodePath := fmt.Sprintf("%s/lock_", z.lockPath)
	ownNode, err := z.conn.Create(nodePath, nil, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	z.ownNodePath = ownNode

	for {
		// 获取所有子节点
		children, _, err := z.conn.Children(z.lockPath)
		if err != nil {
			return err
		}

		// 排序节点
		minNode := ""
		ownNodeName := z.ownNodePath[len(z.lockPath)+1:]
		for _, child := range children {
			if minNode == "" || child < minNode {
				minNode = child
			}
		}

		// 如果自己是最小的节点，则获得锁
		if ownNodeName == minNode {
			glog.Info("获取到锁节点: ", z.ownNodePath)
			return nil
		}

		// 否则监听前一个节点
		predecessor := ""
		for _, child := range children {
			if child < ownNodeName && (predecessor == "" || child > predecessor) {
				predecessor = child
			}
		}

		predecessorPath := fmt.Sprintf("%s/%s", z.lockPath, predecessor)
		_, _, eventCh, err := z.conn.GetW(predecessorPath)
		if err != nil {
			if errors.Is(err, zk.ErrNoNode) {
				// 前驱节点已被删除，重试
				continue
			}
			return err
		}

		// 阻塞等待前驱节点删除事件
		<-eventCh
	}
}

// AcquireNotWatch 获取锁，获取不到时就结束，不监听前驱节点
func (z *ZkLock) AcquireNotWatch() (bool, error) {
	// 确保锁目录存在
	exists, _, err := z.conn.Exists(z.lockPath)
	if err != nil {
		return false, err
	}
	if !exists {
		_, err = z.conn.Create(z.lockPath, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return false, err
		}
	}

	// 创建临时有序节点
	nodePath := fmt.Sprintf("%s/lock_", z.lockPath)
	ownNode, err := z.conn.Create(nodePath, nil, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		return false, err
	}
	z.ownNodePath = ownNode

	for {
		// 获取所有子节点
		children, _, err := z.conn.Children(z.lockPath)
		if err != nil {
			return false, err
		}

		// 排序节点
		minNode := ""
		ownNodeName := z.ownNodePath[len(z.lockPath)+1:]
		for _, child := range children {
			if minNode == "" || child < minNode {
				minNode = child
			}
		}

		// 如果自己是最小的节点，则获得锁
		if ownNodeName == minNode {
			glog.Info("获取到锁节点: ", z.ownNodePath)
			return true, nil
		} else {
			return false, nil
		}
	}
}

// Release 释放锁
func (z *ZkLock) Release() error {
	err := z.conn.Delete(z.ownNodePath, -1)
	if err == nil {
		glog.Info("成功释放锁: ", z.ownNodePath)
	}
	return err
}
