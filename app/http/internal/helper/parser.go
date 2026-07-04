package helper

import (
	"context"
	"log/slog"
	"server-api/app/http/internal/helper/apptypes"

	"github.com/gomooth/pkg/http/httpcontext"

	"github.com/gomooth/xerror"

	"github.com/gin-gonic/gin"
)

// ParseUser 从上下文中解析授权用户
func ParseUser(ctx context.Context) (*httpcontext.User, error) {
	htx, err := httpcontext.Parse(ctx)
	if nil != err {
		return nil, err
	}

	return htx.User(), nil
}

// UserFromContext 从上下文中获取授权用户，解析失败返回零值
func UserFromContext(ctx context.Context) *httpcontext.User {
	user, err := ParseUser(ctx)
	if nil != err {
		slog.Warn("UserFromContext: parse failed", "err", err)
		return &httpcontext.User{}
	}

	return user
}

// ParseAPPHeader 从上下文中解析 APP Header
func ParseAPPHeader(ctx context.Context) (*apptypes.APPHeader, error) {
	gtx, ok := ctx.(*gin.Context)
	if !ok {
		return nil, xerror.New("parse gtx failed")
	}

	var h apptypes.APPHeader
	// 优先从 http header 解析
	if err := gtx.ShouldBindHeader(&h); nil != err {
		return nil, xerror.Wrap(err, "parse global header failed from http header")
	}

	//if h == nil {
	//	// 兼容：如果从 header 获取失败，则从 query 中获取
	//	if err := gtx.ShouldBindQuery(h); nil != err {
	//		return nil, xerror.Wrap(err, "parse global header failed from query string")
	//	}
	//}

	return &h, nil
}

// APPHeaderFromContext 从上下文中获取 APP Header，解析失败返回 nil
func APPHeaderFromContext(ctx context.Context) *apptypes.APPHeader {
	h, err := ParseAPPHeader(ctx)
	if err != nil {
		slog.Warn("APPHeaderFromContext: parse failed", "err", err)
		return nil
	}

	return h
}
