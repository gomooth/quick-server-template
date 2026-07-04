package route

import (
	"server-api/app/openapi/internal/api/system"
	mw "server-api/app/openapi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerSystem(router *gin.Engine) {
	api := new(system.Controller)

	// 需要完整认证的路由
	ra := router.Group("/systems",
		mw.SignDebug(), // 签名调试中间件（在 Auth 之前）
		mw.Auth(),
	)
	{
		ra.GET("/current-times", api.CurrentTime)
	}

	// 仅需 AppID 验证的路由（不需要签名校验）
	rs := router.Group("/systems",
		mw.SignDebug(), // 签名调试中间件
	)
	{
		rs.POST("/signatures", mw.AuthWithoutSign(), api.Sign)
	}
}
