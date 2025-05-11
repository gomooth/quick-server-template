package helper

import (
	"context"
	"server-api/repository/platform"

	"github.com/gin-gonic/gin"

	"github.com/save95/xerror"
)

var OpenAPPIDKey = "__open_app"

func MustParseOpenAPP(ctx context.Context) (*platform.OpenAPP, error) {
	gtx, ok := ctx.(*gin.Context)
	if !ok {
		return nil, xerror.New("to gtx failed")
	}

	v, ok := gtx.Get(OpenAPPIDKey)
	if !ok {
		return nil, nil
	}

	return v.(*platform.OpenAPP), nil
}
