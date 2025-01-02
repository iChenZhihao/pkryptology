package initial

import "github.com/coinbase/kryptology/service/zkp"

func Run() {
	// 加载配置信息
	LoadConfig()

	// 初始化Zk连接
	InitZookeeper()

	// 将服务注册到Zk中
	zkp.Register()

	// 启动Gin
	Router()
}
