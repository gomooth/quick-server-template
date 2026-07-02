package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomooth/pkg/http/restful"
)

type Controller struct{}

// Ping 兼容旧 /ping 端点，返回 pong
func (c *Controller) Ping(ctx *gin.Context) {
	res := new(service).ping()
	rru := restful.NewResponse(ctx)
	rru.WithBody(res.Message)
}

// Liveness 存活探针，进程存活即返回 200
func (c *Controller) Liveness(ctx *gin.Context) {
	res := new(service).liveness()
	rru := restful.NewResponse(ctx)
	rru.Retrieve(res)
}

// Readiness 就绪探针，依赖全部可用才返回 200
func (c *Controller) Readiness(ctx *gin.Context) {
	res := new(service).readiness(ctx.Request.Context())

	if res.Status == healthFail {
		ctx.JSON(http.StatusServiceUnavailable, res)
		return
	}

	rru := restful.NewResponse(ctx)
	rru.Retrieve(res)
}

// Endpoint 接受数据
func (c *Controller) Endpoint(ctx *gin.Context) {
	rru := restful.NewResponse(ctx)
	rru.WithMessage("success")
}
