package health

import (
	healthsvc "server-api/service/health"

	"github.com/gin-gonic/gin"
	"github.com/gomooth/pkg/http/restful"
)

type Controller struct{}

// Ping 兼容旧 /ping 端点
func (c *Controller) Ping(ctx *gin.Context) {
	res := new(healthsvc.Service).Ping()
	rru := restful.NewResponse(ctx)
	rru.WithBody(res.Message)
}

// Endpoint 接受数据
func (c *Controller) Endpoint(ctx *gin.Context) {
	rru := restful.NewResponse(ctx)
	rru.WithMessage("success")
}
