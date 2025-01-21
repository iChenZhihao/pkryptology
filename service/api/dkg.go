package api

import (
	"encoding/json"
	"github.com/coinbase/kryptology/service/gg20/node"
	"github.com/coinbase/kryptology/service/respvo"
	"github.com/gin-gonic/gin"
	"sync"
)

var dkgController *DkgController
var once sync.Once

type DkgController struct {
	operator *node.DkgOperator
}

func GetDkgController() *DkgController {
	once.Do(func() {
		dkgController = &DkgController{}
		dkgController.operator = node.GetDkgOperator()
	})
	return dkgController
}

func (d *DkgController) DoRound1(c *gin.Context) {
	//glog.Info("DkgRound1触发接口被调用了")
	d.operator.TriggeredToStartRound1()
	respvo.Success("", nil, c)
}

func (d *DkgController) DoRound1Recv(c *gin.Context) {
	//glog.Info("DkgRound1Recv接收接口被调用了")
	data, err := c.GetRawData()
	if err != nil {
		respvo.Failed("获取DkgRound1请求体失败", c)
		return
	}
	recvBody := &node.DkgRound1Recv{}
	err = json.Unmarshal(data, recvBody)
	if err != nil {
		respvo.Failed("反序列化DkgRound1失败", c)
		return
	}
	d.operator.RecvDkgRound1(*recvBody)
	respvo.Success("", nil, c)
}

func (d *DkgController) DoRound2Recv(c *gin.Context) {
	//glog.Info("DkgRound2Recv接收接口被调用了")
	data, err := c.GetRawData()
	if err != nil {
		respvo.Failed("获取DkgRound2请求体失败", c)
		return
	}
	recvBody := &node.DkgRound2Recv{}
	err = json.Unmarshal(data, recvBody)
	if err != nil {
		respvo.Failed("反序列化DkgRound2失败", c)
		return
	}
	d.operator.RecvDkgRound2(*recvBody)
	respvo.Success("", nil, c)
}

func (d *DkgController) DoRound3Recv(c *gin.Context) {
	//glog.Info("DkgRound2Recv接收接口被调用了")
	data, err := c.GetRawData()
	if err != nil {
		respvo.Failed("获取DkgRound3请求体失败", c)
		return
	}
	recvBody := &node.DkgRound3Recv{}
	err = json.Unmarshal(data, recvBody)
	if err != nil {
		respvo.Failed("反序列化DkgRound3失败", c)
		return
	}
	d.operator.RecvDkgRound3(*recvBody)
	respvo.Success("", nil, c)
}
