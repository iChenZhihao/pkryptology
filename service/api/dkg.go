package api

import (
	"encoding/json"
	"github.com/coinbase/kryptology/service/gg20/dkg"
	"github.com/coinbase/kryptology/service/response"
	"github.com/gin-gonic/gin"
	//"github.com/golang/glog"
	"sync"
	//"github.com/golang/glog"
)

var dkgController *DkgController
var once sync.Once

type DkgController struct {
	operator *dkg.DkgOperator
}

func GetDkgController() *DkgController {
	once.Do(func() {
		dkgController = &DkgController{}
		dkgController.operator = dkg.GetDkgOperator()
	})
	return dkgController
}

func (d *DkgController) DoRound1(c *gin.Context) {
	//glog.Info("DkgRound1触发接口被调用了")
	value := c.Query("secret")
	secret := dkg.Base64DecodeSecret(value)
	d.operator.TriggeredToStartRound1(secret)
	//if err != nil {
	//	glog.Errorf("被调用启动Round1失败：%v\n", err.Error())
	//	response.Failed("", c)
	//	return
	//}
	response.Success("", nil, c)
}

func (d *DkgController) DoRound1Recv(c *gin.Context) {
	//glog.Info("DkgRound1Recv接收接口被调用了")
	data, err := c.GetRawData()
	if err != nil {
		response.Failed("获取DkgRound1请求体失败", c)
		return
	}
	recvBody := &dkg.DkgRound1Recv{}
	err = json.Unmarshal(data, recvBody)
	if err != nil {
		response.Failed("反序列化DkgRound1失败", c)
		return
	}
	d.operator.RecvDkgRound1(*recvBody)
	response.Success("", nil, c)
}

func (d *DkgController) DoRound2Recv(c *gin.Context) {
	//glog.Info("DkgRound2Recv接收接口被调用了")
	data, err := c.GetRawData()
	if err != nil {
		response.Failed("获取DkgRound2请求体失败", c)
		return
	}
	recvBody := &dkg.DkgRound2Recv{}
	err = json.Unmarshal(data, recvBody)
	if err != nil {
		response.Failed("反序列化DkgRound2失败", c)
		return
	}
	d.operator.RecvDkgRound2(*recvBody)
	response.Success("", nil, c)
}
