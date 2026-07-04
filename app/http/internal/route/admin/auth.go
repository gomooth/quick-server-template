package admin

import (
	"server-api/global"
	"server-api/app/http/internal/api/admin/auth"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/pkg/http/jwt/jwtstore"
	"github.com/gomooth/pkg/http/middleware"
)

func registerAuth(router *gin.Engine) {
	api := auth.Controller{}

	// 公开：登录
	ra := router.Group(
		"/admin/auth",
		middleware.RESTFul(global.ApiVersionLatest),
	)
	{
		ra.POST("/tokens", api.Token)
	}

	// 认证：登出 + 改密
	ra2 := router.Group(
		"/admin/auth",
		middleware.RESTFul(global.ApiVersionLatest),
		middleware.JWTStatefulWith(
			[]byte(global.Config.App.Secret),
			global.NewRole,
			jwtstore.NewMultiRedisStore(global.SessionStoreClient),
		),
		middleware.WithRole(global.RoleSuper),
	)
	{
		ra2.DELETE("/tokens", api.Logout)
		ra2.PUT("/passwords", api.ChangePwd)
	}
}
