package helper

import (
	"context"
	"server-api/app/http/internal/helper/apptypes"
	"server-api/global"

	"github.com/gomooth/pkg/http/httpcontext"

	"github.com/save95/xerror"

	"github.com/gin-gonic/gin"
)

// MustParseUser 从上下文中解析授权用户，否则报错
func MustParseUser(ctx context.Context) (*httpcontext.User, error) {
	htx, err := httpcontext.MustParse(ctx)
	if nil != err {
		return nil, err
	}

	return htx.User(), nil
}

// ParseUser 从上下文中解析授权用户
func ParseUser(ctx context.Context) *httpcontext.User {
	user, err := MustParseUser(ctx)
	if nil != err {
		global.Log.Warningf("ParseUser: parse failed, err=%+v", err)
		return &httpcontext.User{}
	}

	return user
}

func MustParseAPPHeader(ctx context.Context) (*apptypes.APPHeader, error) {
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

func ParseAPPHeader(ctx context.Context) *apptypes.APPHeader {
	h, err := MustParseAPPHeader(ctx)
	if err != nil {
		global.Log.Warningf("ParseAppHeader: parse failed, err=%+v", err)
		return nil
	}

	return h
}
