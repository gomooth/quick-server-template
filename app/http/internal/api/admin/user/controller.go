package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gomooth/pkg/http/restful"
	"github.com/gomooth/utils/valutil"
	"github.com/gomooth/xerror"
	"github.com/gomooth/xerror/xcode"
)

type Controller struct{}

func (c *Controller) Paginate(ctx *gin.Context) {
	rru := restful.NewResponse(ctx)

	var in paginateRequest
	if err := ctx.ShouldBindQuery(&in); nil != err {
		rru.WithError(err)
		return
	}

	records, total, err := new(service).Paginate(ctx, &in)
	if nil != err {
		rru.WithError(err)
		return
	}

	rru.ListWithPagination(total, records)
}

func (c *Controller) Create(ctx *gin.Context) {
	rru := restful.NewResponse(ctx)

	var in createRequest
	if err := ctx.ShouldBindJSON(&in); nil != err {
		rru.WithError(err)
		return
	}

	if err := in.Validate(); nil != err {
		rru.WithError(xerror.NewXCode(xcode.RequestParamError, err.Error()))
		return
	}

	record, err := new(service).Create(ctx, &in)
	if nil != err {
		rru.WithError(err)
		return
	}

	rru.Post(record)
}

func (c *Controller) Modify(ctx *gin.Context) {
	rru := restful.NewResponse(ctx)

	var in modifyRequest
	if err := ctx.ShouldBindJSON(&in); nil != err {
		rru.WithError(err)
		return
	}

	if err := in.Validate(); nil != err {
		rru.WithError(xerror.NewXCode(xcode.RequestParamError, err.Error()))
		return
	}

	id := valutil.Int(ctx.Param("id"))
	record, err := new(service).Modify(ctx, uint(id), &in)
	if nil != err {
		rru.WithError(err)
		return
	}

	rru.Post(record)
}
