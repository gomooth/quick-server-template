package auth

import (
	"server-api/service/lang"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/pkg/http/jwt"
	"github.com/gomooth/pkg/http/restful"
)

type Controller struct {
}

func (c *Controller) Token(ctx *gin.Context) {
	rru := restful.NewResponse(
		ctx,
		restful.WithResponseErrorMsgHandler(lang.Handler()),
	)

	var in createTokenRequest
	if err := ctx.ShouldBindJSON(&in); nil != err {
		rru.WithError(err)
		return
	}

	token, err := new(service).Login(ctx, &in)
	if err != nil {
		rru.WithError(err)
		return
	}

	rru.SetHeader(jwt.TokenHeaderKey, token.AccessToken)
	rru.Post(token)
}

func (c *Controller) Logout(ctx *gin.Context) {
	rru := restful.NewResponse(
		ctx,
		restful.WithResponseErrorMsgHandler(lang.Handler()),
	)

	err := new(service).Logout(ctx)
	rru.Delete(err)
}

func (c *Controller) ChangePwd(ctx *gin.Context) {
	rru := restful.NewResponse(
		ctx,
		restful.WithResponseErrorMsgHandler(lang.Handler()),
	)

	var in changePwdRequest
	if err := ctx.ShouldBindJSON(&in); nil != err {
		rru.WithError(err)
		return
	}

	if err := new(service).ChangePwd(ctx, &in); nil != err {
		rru.WithError(err)
		return
	}

	rru.WithMessage("success")
}
