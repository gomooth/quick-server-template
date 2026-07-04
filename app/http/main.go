package http

import (
	"server-api/global"
	"server-api/app/http/internal/api/ping"
	"server-api/app/http/internal/route"
	"server-api/app/http/internal/route/admin"
	"server-api/app/http/internal/route/user"

	"github.com/gomooth/pkg/http/middleware"

	"github.com/gin-gonic/gin"
)

// RouteRegister 注册所有的路由
// 这里，路由请按模块分开写，一个模块一个文件（建议一个模块中分多个函数来写子模块）
// 单个模块的路由注册使用私有方法，不对外暴露
func RouteRegister(router *gin.Engine) {
	// 静态文件
	router.Static("/storage", "storage/public")

	router.Any("/ping", new(ping.Controller).Ping)
	router.Any("/endpoint", middleware.HttpPrinter(global.Log), new(ping.Controller).Endpoint)

	// 注册路由
	admin.Register(router)
	user.Register(router)
	route.RegisterFile(router)
}

// Release 释放资源
func Release() error {

	return nil
}

