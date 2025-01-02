package initial

import (
	"fmt"
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

	_ = engine.Run(fmt.Sprintf(":%s", global.Config.Server.Port))
}
