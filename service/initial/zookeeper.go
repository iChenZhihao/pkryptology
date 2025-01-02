package initial

import (
	"github.com/coinbase/kryptology/service/global"
	"github.com/golang/glog"
	"time"
)

func InitZookeeper() {
	zkcfg := global.Config.Zookeeper
	glog.Info("Loaded config: %+v\n", zkcfg)

	zkManager := global.GetZkManager()
	err := zkManager.Init(zkcfg.Servers, time.Duration(zkcfg.SessionTimeout)*time.Millisecond)
	if err != nil {
		glog.Errorf("Failed to initialize Zookeeper connection: %v", err)
	}

	glog.Info("Connected to Zookeeper!")

	//if exists, _, err := conn.Exists("/gg20/nodes"); err != nil {
	//	glog.Error(err.Error())
	//} else if !exists {
	//	if _, err := conn.Create("/gg20/nodes", nil, 0, zk.WorldACL(zk.PermAll)); err != nil {
	//		glog.Error(err.Error())
	//	}
	//}

}
