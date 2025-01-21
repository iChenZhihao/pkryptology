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
	_, err = c.signer.SignMsg([]byte(toSign.Message))
	if err != nil {
		respvo.Failed("签名失败: "+err.Error(), ctx)
		return
	}
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
