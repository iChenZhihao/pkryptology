package initial

import (
	"fmt"
	"github.com/coinbase/kryptology/service/api"
	"github.com/coinbase/kryptology/service/component"
	"github.com/coinbase/kryptology/service/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Router() {
	engine := gin.Default()

	// 开启跨域
	engine.Use(component.Cors())

	engine.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404 not found")
	})

	dkgGroup := engine.Group("/dkg")

	{
		dkgGroup.GET("/round1", api.GetDkgController().DoRound1)
		dkgGroup.POST("/round1/recv", api.GetDkgController().DoRound1Recv)
		dkgGroup.POST("/round2/recv", api.GetDkgController().DoRound2Recv)
		dkgGroup.POST("/round3/recv", api.GetDkgController().DoRound3Recv)
	}

	signGroup := engine.Group("/sign")
	{
		signGroup.POST("/signMsg", api.GetSignController().SignMsg)
		signGroup.POST("/candidate", api.GetSignController().RecvAskCandidateInfo)
		signGroup.POST("/round1", api.GetSignController().DoSignRound1)
		signGroup.POST("/round1/recv", api.GetSignController().DoSignRound1Recv)
		signGroup.POST("/round2/recv", api.GetSignController().DoSignRound2Recv)
		signGroup.POST("/round3/recv", api.GetSignController().DoSignRound3Recv)
		signGroup.POST("/round4/recv", api.GetSignController().DoSignRound4Recv)
		signGroup.POST("/round5/recv", api.GetSignController().DoSignRound5Recv)
		signGroup.POST("/round6/recv", api.GetSignController().DoSignRound6Recv)
	}

	_ = engine.Run(fmt.Sprintf(":%s", global.Config.Server.Port))
}
