package middleware

import (
	"server-api/app/openapi/internal/helper"
	"server-api/global"

	"github.com/gin-gonic/gin"
	"github.com/gomooth/pkg/http/restful"
)

// SignDebug 签名调试中间件，挂载在 Auth 之前
// 当配置启用且请求 AppID 匹配调试 AppID 时，用真实 path 计算签名并返回，然后 Abort
func SignDebug() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !global.Config.Server.OpenAPI.SignDebugEnabled {
			c.Next()
			return
		}

		h, err := helper.ParseHeader(c)
		if err != nil {
			c.Next() // Header 解析失败，交给后续 auth 处理
			return
		}

		if h.AppID != global.Config.Server.OpenAPI.SignDebugAppID {
			c.Next() // 非调试 AppID，跳过
			return
		}

		// 匹配调试 AppID → 用真实 path 计算签名
		qs, body := helper.ExtractRequestParams(c)
		result := helper.Sign(h.AppID, global.Config.Server.OpenAPI.SignDebugSecret,
			h.Timestamp, c.Request.URL.Path, qs, body)
		result.Input = h.Sign
		result.Success = result.Input == result.Signature

		rru := restful.NewResponse(c)
		rru.Retrieve(result)
		c.Abort() // 不继续走 auth 和业务逻辑
	}
}
