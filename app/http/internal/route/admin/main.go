package admin

import (
	"server-api/global"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/pkg/http/jwt/jwtstore"
	"github.com/gomooth/pkg/http/middleware"
)

func Register(router *gin.Engine) {
	// 公开路由：admin 登录
	registerAuth(router)

	// 受保护路由：JWT + RoleSuper
	ra := router.Group(
		"/admin",
		middleware.RESTFul(global.ApiVersionLatest),
		middleware.JWTStatefulWith(
			[]byte(global.Config.App.Secret),
			global.NewRole,
			jwtstore.NewMultiRedisStore(global.SessionStoreClient), // 多地登录
		),
		middleware.WithRole(global.RoleSuper),
	)

	registerUser(ra)
}
