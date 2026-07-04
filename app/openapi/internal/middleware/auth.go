package middleware

import (
	"log/slog"
	"math"
	"server-api/app/openapi/internal/helper"
	"server-api/app/openapi/internal/helper/ecode"
	"server-api/global"
	"server-api/repository/platform/pattr"
	"server-api/repository/platform/pcache"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/pkg/http/restful"
	"github.com/gomooth/utils/valutil"

	"github.com/gomooth/xerror"
	"github.com/gomooth/xerror/xcode"
)

func Auth() gin.HandlerFunc {
	return auth(false)
}

func AuthWithoutSign() gin.HandlerFunc {
	return auth(true)
}

func auth(withoutSign bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		rru := restful.NewResponse(c)

		h, err := helper.ParseHeader(c)
		if err != nil {
			rru.WithError(err)
			return
		}
		if err := h.Validate(); nil != err {
			rru.WithError(err)
			return
		}

		// 检查时间是否超过 5分钟
		if !withoutSign {
			requestAt := time.Unix(int64(valutil.Int(h.Timestamp)), 0)
			if math.Abs(time.Now().Sub(requestAt).Minutes()) > 5 {
				rru.WithError(xerror.NewXCode(ecode.RequestExpired))
				return
			}
		}

		// 获取应用
		app, err := pcache.NewOpenAPP().FirstByAppID(c, h.AppID)
		if err != nil {
			if xerror.IsXCode(err, xcode.DBRecordNotFound) {
				rru.WithError(xerror.NewXCode(ecode.OpenAPIClosed))
				return
			}

			rru.WithError(xerror.WrapWithXCode(err, ecode.InternalError))
			return
		}
		if app.State != pattr.OpenAPPStateNormal {
			rru.WithError(xerror.NewXCode(ecode.OpenAPIClosed))
			return
		}

		// 写入自定义上下文
		c.Set(helper.OpenAPPIDKey, app)

		if withoutSign {
			c.Next()
			return
		}

		// 检查签名
		qs, body := helper.ExtractRequestParams(c)
		path := c.Request.URL.Path
		sign := helper.Sign(h.AppID, app.AppSecret, h.Timestamp, path, qs, body)
		if h.Sign != sign.Signature {
			if !global.Env().IsProd() {
				slog.Debug("sign failed", "timestamp", h.Timestamp, "input", h.Sign, "sign", sign.Signature)
			}
			rru.WithError(xerror.NewXCode(ecode.SignError))
			return
		}

		c.Next()
	}
}
