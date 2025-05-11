package route

import (
	"server-api/app/http/internal/api/auth"
	"server-api/app/http/internal/helper"
	"server-api/global"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/pkg/http/jwt/jwtstore"
	"github.com/gomooth/pkg/http/middleware"
)

// RegisterAuth 注册鉴权路由
func RegisterAuth(router *gin.Engine) {
	api := auth.Controller{}

	ra := router.Group(
		"/auth",
		middleware.RESTFul(global.ApiVersionLatest),
	)
	{
		// 创建 Token
		ra.POST("/tokens", api.Token)
	}

	ra2 := router.Group(
		"/auth",
		middleware.RESTFul(global.ApiVersionLatest),
		middleware.JWTStatefulWith(
			helper.JWTOption(false),
			jwtstore.NewMultiRedisStore(global.SessionStoreClient), // 多地登录
		),
		middleware.WithRole(global.RoleUser),
	)
	{
		// 退出登陆
		ra2.DELETE("/tokens", api.Logout)
		// 修改密码
		ra2.PUT("/passwords", api.ChangePwd)
	}
}
