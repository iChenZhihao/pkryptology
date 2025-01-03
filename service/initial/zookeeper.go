package initial

import (
	"github.com/coinbase/kryptology/service/global"
	"github.com/coinbase/kryptology/service/zkp"
	"github.com/golang/glog"
	"time"
)

func InitZookeeper() {
	zkcfg := global.Config.Zookeeper
	glog.Info("Loaded config: %+v\n", zkcfg)

	zkManager := zkp.GetZkManager()
	err := zkManager.Init(zkcfg.Servers, time.Duration(zkcfg.SessionTimeout)*time.Millisecond)
	if err != nil {
		glog.Errorf("Failed to initialize Zookeeper connection: %v", err)
	}

	glog.Info("Connected to Zookeeper!")

}
