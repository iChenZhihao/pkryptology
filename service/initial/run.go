package initial

import (
	"github.com/coinbase/kryptology/service/zkp"
	"github.com/golang/glog"
)

func Run() {
	// 加载配置信息
	//LoadConfig()

	// 初始化Zk连接
	InitZookeeper()

	// 将服务注册到Zk中
	err := zkp.GetZkManager().Register()
	if err != nil {
		glog.Errorf("将服务注册到Zk中失败\n")
		return
	}
	go zkp.GetZkManager().MonitorNodeChanges()
	defer zkp.GetZkManager().Close()

	// 启动Gin
	Router()
}
