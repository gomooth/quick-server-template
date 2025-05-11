package openapi

import (
	"server-api/app/openapi/internal/api/ping"
	"server-api/app/openapi/internal/route"
	"server-api/global"

	"github.com/gomooth/pkg/http/middleware"

	"github.com/gin-gonic/gin"
)

// Register 注册所有的路由
// 这里，路由请按模块分开写，一个模块一个文件（建议一个模块中分多个函数来写子模块）
// 单个模块的路由注册使用私有方法，不对外暴露
func Register(router *gin.Engine) {
	router.Any("/ping", ping.Controller{}.Ping)
	router.Any("/endpoint", middleware.HttpPrinter(global.Log), ping.Controller{}.Endpoint)

	// 注册路由
	route.RegisterSystem(router)
}
