package admin

import (
	"server-api/app/http/internal/api/admin/user"

	"github.com/gin-gonic/gin"
)

func registerUser(router gin.IRouter) {
	api := user.Controller{}
	ra := router.Group("/users")
	{
		ra.GET("", api.Paginate)
		ra.POST("", api.Create)
		ra.PUT("/:id", api.Modify)
	}
}
