package api

import (
	"encoding/json"
	"github.com/coinbase/kryptology/service/gg20/node"
	"github.com/coinbase/kryptology/service/respvo"
	"github.com/gin-gonic/gin"
	"sync"
)

var signerOnce sync.Once
var signController *SignController

type SignController struct {
	signer *node.SignOperator
}

type ToSignMsg struct {
	Message string `json:"message"`
}

func GetSignController() *SignController {
	signerOnce.Do(func() {
		signController = &SignController{}
		signController.signer = node.GetSignOperator()
	})
	return signController
}

func (c *SignController) SignMsg(ctx *gin.Context) {
	data, err := ctx.GetRawData()
	if err != nil {
		respvo.Failed("获取待签名信息失败", ctx)
		return
	}
	toSign := &ToSignMsg{}
	err = json.Unmarshal(data, toSign)
	if err != nil {
		respvo.Failed("数据格式异常，反序列化失败", ctx)
		return
	}
	signature, err := c.signer.SignMsg([]byte(toSign.Message))
	if err != nil {
		respvo.Failed("签名失败: "+err.Error(), ctx)
		return
	}
	respvo.Success("", signature, ctx)
}

func (c *SignController) RecvAskCandidateInfo(ctx *gin.Context) {
	data, err := ctx.GetRawData()
	if err != nil {
		respvo.Failed("获取WorkId信息失败", ctx)
		return
	}
	workId := string(data)
	info := c.signer.RecvAskCandidateInfo(workId)
	respvo.Success("", info, ctx)
}

func (c *SignController) DoSignRound1(ctx *gin.Context) {
	data, err := ctx.GetRawData()
	if err != nil {
		respvo.Failed("获取请求体信息失败", ctx)
		return
	}
	info := &node.StartSignInfo{}
	err = json.Unmarshal(data, info)
	if err != nil {
		respvo.Failed("数据格式异常，反序列化失败", ctx)
		return
	}
	err = c.signer.TriggeredToStartSign(info.WorkId, info.Cosigner, info.HashMsg)
	if err != nil {
		respvo.Failed("触发签名失败: "+err.Error(), ctx)
		return
	}
	respvo.Success("", nil, ctx)
}

func (c *SignController) DoSignRound1Recv(ctx *gin.Context) {
	data, err := ctx.GetRawData()
	if err != nil {
		respvo.Failed("获取请求体信息失败", ctx)
		return
	}
	info := &node.SignRound1Recv{}
	err = json.Unmarshal(data, info)
	if err != nil {
		respvo.Failed("数据格式异常，反序列化失败", ctx)
		return
	}
	err = c.signer.DoSignRound1Recv(info)
	if err != nil {
		respvo.Failed("接收广播信息失败: "+err.Error(), ctx)
		return
	}
	respvo.Success("", nil, ctx)
}

func (c *SignController) DoSignRound2Recv(ctx *gin.Context) {
	data, err := ctx.GetRawData()
	if err != nil {
		respvo.Failed("获取请求体信息失败", ctx)
		return
	}
	info := &node.SignRound2Recv{}
	err = json.Unmarshal(data, info)
	if err != nil {
		respvo.Failed("数据格式异常，反序列化失败", ctx)
		return
	}
	err = c.signer.DoSignRound2Recv(info)
	if err != nil {
		respvo.Failed("接收广播信息失败: "+err.Error(), ctx)
		return
	}
	respvo.Success("", nil, ctx)
}

func (c *SignController) DoSignRound3Recv(ctx *gin.Context) {
	data, err := ctx.GetRawData()
	if err != nil {
		respvo.Failed("获取请求体信息失败", ctx)
		return
	}
	info := &node.SignRound3Recv{}
	err = json.Unmarshal(data, info)
	if err != nil {
		respvo.Failed("数据格式异常，反序列化失败", ctx)
		return
	}
	err = c.signer.DoSignRound3Recv(info)
	if err != nil {
		respvo.Failed("接收广播信息失败: "+err.Error(), ctx)
		return
	}
	respvo.Success("", nil, ctx)
}

func (c *SignController) DoSignRound4Recv(ctx *gin.Context) {
	data, err := ctx.GetRawData()
	if err != nil {
		respvo.Failed("获取请求体信息失败", ctx)
		return
	}
	info := &node.SignRound4Recv{}
	err = json.Unmarshal(data, info)
	if err != nil {
		respvo.Failed("数据格式异常，反序列化失败", ctx)
		return
	}
	err = c.signer.DoSignRound4Recv(info)
	if err != nil {
		respvo.Failed("接收广播信息失败: "+err.Error(), ctx)
		return
	}
	respvo.Success("", nil, ctx)
}

func (c *SignController) DoSignRound5Recv(ctx *gin.Context) {
	data, err := ctx.GetRawData()
	if err != nil {
		respvo.Failed("获取请求体信息失败", ctx)
		return
	}
	info := &node.SignRound5Recv{}
	err = json.Unmarshal(data, info)
	if err != nil {
		respvo.Failed("数据格式异常，反序列化失败", ctx)
		return
	}
	err = c.signer.DoSignRound5Recv(info)
	if err != nil {
		respvo.Failed("接收广播信息失败: "+err.Error(), ctx)
		return
	}
	respvo.Success("", nil, ctx)
}

func (c *SignController) DoSignRound6Recv(ctx *gin.Context) {
	data, err := ctx.GetRawData()
	if err != nil {
		respvo.Failed("获取请求体信息失败", ctx)
		return
	}
	info := &node.SignRound6Recv{}
	err = json.Unmarshal(data, info)
	if err != nil {
		respvo.Failed("数据格式异常，反序列化失败", ctx)
		return
	}
	err = c.signer.DoSignRound6Recv(info)
	if err != nil {
		respvo.Failed("接收广播信息失败: "+err.Error(), ctx)
		return
	}
	respvo.Success("", nil, ctx)
}
