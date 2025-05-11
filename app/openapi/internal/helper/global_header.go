package helper

import (
	"context"
	"server-api/app/openapi/internal/helper/ecode"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/save95/xerror"
)

type Header struct {
	AppID      string `header:"X-App-Id" form:"appId"`
	AppVersion string `header:"X-Api-Version" form:"appVersion"`
	SignType   string `header:"X-Sign-Type" form:"signType"`
	Timestamp  string `header:"X-Timestamp" form:"timestamp"`
	Sign       string `header:"X-Sign" form:"sign"`
}

func (h *Header) Validate() error {
	if len(h.AppID) == 0 || len(h.AppVersion) == 0 || len(h.SignType) == 0 ||
		len(h.Timestamp) == 0 || len(h.Sign) == 0 {
		return xerror.WithXCode(ecode.MissRequired)
	}
	if h.AppVersion != Version {
		return xerror.WithXCode(ecode.RequiredParamError)
	}
	if h.SignType != "sha1" {
		return xerror.WithXCode(ecode.RequiredParamError)
	}
	if len(h.Timestamp) != 10 {
		return xerror.WithXCode(ecode.RequiredParamError)
	}
	if _, err := strconv.ParseInt(h.Timestamp, 10, 64); err != nil {
		return xerror.WithXCode(ecode.RequiredParamError)
	}
	return nil
}

func MustParseHeader(ctx context.Context) (*Header, error) {
	gtx, ok := ctx.(*gin.Context)
	if !ok {
		return nil, xerror.New("parse gtx failed")
	}

	var h Header
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
