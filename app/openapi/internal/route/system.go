package route

import (
	"server-api/app/openapi/internal/api/system"
	mw "server-api/app/openapi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterSystem(router *gin.Engine) {
	ra := router.Group(
		"/systems",
	)

	api := system.Controller{}
	ra.GET("/current-times", mw.Auth(), api.CurrentTime)
	ra.POST("/signatures", mw.AuthWithoutSign(), api.Sign)
}
