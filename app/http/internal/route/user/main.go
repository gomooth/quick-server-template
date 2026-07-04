package user

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	// 公开路由：user 登录
	registerAuth(router)

	// 受保护路由（JWT + RoleUser）
	// 未来 user 模块在此注册
	// g := router.Group("/user",
	//     middleware.JWTStatefulWith(...),
	//     middleware.WithRole(global.RoleUser),
	// )
	// registerXxx(g)
}
