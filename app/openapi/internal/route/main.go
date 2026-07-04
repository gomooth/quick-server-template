package route

import "github.com/gin-gonic/gin"

// Register 注册所有路由
func Register(router *gin.Engine) {
	registerHealth(router)
	registerSystem(router)
}
