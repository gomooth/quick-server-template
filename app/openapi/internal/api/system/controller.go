package system

import (
	"server-api/app/openapi/internal/helper"

	"github.com/gin-gonic/gin"
	"github.com/gomooth/pkg/http/restful"
)

type Controller struct{}

func (c *Controller) CurrentTime(ctx *gin.Context) {
	rru := restful.NewResponse(ctx)
	res, err := new(service).CurrentTime(ctx)
	if err != nil {
		rru.WithError(err)
		return
	}
	rru.Retrieve(res)
}

func (c *Controller) Sign(ctx *gin.Context) {
	rru := restful.NewResponse(ctx)

	h, err := helper.ParseHeader(ctx)
	if err != nil {
		rru.WithError(err)
		return
	}

	res, err := new(service).Sign(ctx, h)
	if err != nil {
		rru.WithError(err)
		return
	}
	rru.Retrieve(res)
}
