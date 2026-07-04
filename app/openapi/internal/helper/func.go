package helper

import (
	"context"
	"server-api/repository/platform/pmodel"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/xerror"
)

var OpenAPPIDKey = "__open_app"

func ParseOpenAPP(ctx context.Context) (*pmodel.OpenAPP, error) {
	gtx, ok := ctx.(*gin.Context)
	if !ok {
		return nil, xerror.New("to gtx failed")
	}

	v, ok := gtx.Get(OpenAPPIDKey)
	if !ok {
		return nil, nil
	}

	return v.(*pmodel.OpenAPP), nil
}
