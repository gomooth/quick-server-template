package route

import (
	"server-api/app/openapi/internal/api/health"
	"server-api/global"

	"github.com/gomooth/pkg/http/middleware"

	"github.com/gin-gonic/gin"
)

func registerHealth(router *gin.Engine) {
	ctrl := new(health.Controller)
	router.GET("/ping", ctrl.Ping)
	router.Any("/endpoint", middleware.HttpPrinter(global.Log), ctrl.Endpoint)
}
