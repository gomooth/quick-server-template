package user

import (
	"server-api/global"
	"server-api/app/http/internal/api/user/auth"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/pkg/http/jwt/jwtstore"
	"github.com/gomooth/pkg/http/middleware"
)

func registerAuth(router *gin.Engine) {
	api := auth.Controller{}

	// 公开：登录
	ra := router.Group(
		"/user/auth",
		middleware.RESTFul(global.ApiVersionLatest),
	)
	{
		ra.POST("/tokens", api.Token)
	}

	// 认证：登出 + 改密
	ra2 := router.Group(
		"/user/auth",
		middleware.RESTFul(global.ApiVersionLatest),
		middleware.JWTStatefulWith(
			[]byte(global.Config.App.Secret),
			global.NewRole,
			jwtstore.NewSingleRedisStore(global.SessionStoreClient), // 单端登录
		),
		middleware.WithRole(global.RoleUser),
	)
	{
		ra2.DELETE("/tokens", api.Logout)
		ra2.PUT("/passwords", api.ChangePwd)
	}
}
