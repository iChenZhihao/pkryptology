package zkp

import (
	"github.com/go-zookeeper/zk"
	"github.com/golang/glog"
	"sync"
	"time"
)

var (
	zkClient *ZkManager
	zkOnce   sync.Once
)

func GetZkManager() *ZkManager {
	zkOnce.Do(func() {
		zkClient = &ZkManager{}
	})
	return zkClient
}

type ZkManager struct {
	conn          *zk.Conn
	mu            sync.Mutex
	lastNodeCount int
}

func (m *ZkManager) Init(servers []string, sessionTimeout time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果连接已经初始化，则直接返回
	if m.conn != nil {
		return nil
	}

	conn, events, err := zk.Connect(servers, sessionTimeout)
	if err != nil {
		return err
	}

	go func() {
		for event := range events {
			glog.Infof("Zookeeper事件: %+v", event)
		}
	}()

	m.conn = conn
	glog.Info("Zk连接已初始化!")
	return nil
}

// GetConn 获取连接对象
func (m *ZkManager) GetConn() *zk.Conn {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.conn
}

// Close 关闭连接
func (m *ZkManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.conn != nil {
		m.conn.Close()
		m.conn = nil
		glog.Info("Zookeeper connection closed")
	}
}
