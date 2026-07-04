package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"math"
	"server-api/global"
	"server-api/app/openapi/internal/helper"
	"server-api/app/openapi/internal/helper/ecode"
	"server-api/repository/platform/pattr"
	"server-api/repository/platform/pcache"
	"strings"
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

var (
	signDebuggerAppID     = "sign-debugger"
	signDebuggerAppSecret = "sign-debugger-secret"
)

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

		// 开启签名调试，直接返回校验结果，不返回数据
		if h.AppID == signDebuggerAppID {
			sign := signDebugger(c, h.Timestamp)
			bs, _ := json.Marshal(sign)
			rru.WithBody(string(bs))
			c.Abort()
			return
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

		path := c.Request.URL.Path

		// 获得请求参数
		qs := make(map[string]string)
		for key := range c.Request.URL.Query() {
			v := c.Request.URL.Query().Get(key)
			if len(v) > 0 {
				qs[key] = v
			}
		}

		bodyRaw, _ := io.ReadAll(c.Request.Body)
		body := strings.TrimSpace(string(bodyRaw)) // 去掉换行

		// 重新写入 request body
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyRaw))

		// 检查签名
		sign := helper.Sign(h.AppID, app.AppSecret, h.Timestamp, path, qs, body)
		if h.Sign != sign.Signature {
			if !global.Env().IsProd() {
				slog.Debug("sign failed", "timestamp", h.Timestamp, "input", h.Sign, "sign", sign.Signature)
			}
			rru.WithError(xerror.NewXCode(ecode.SignError))
			return
		}
	}
}

func signDebugger(c *gin.Context, ts string) *helper.SignResult {
	path := c.Request.URL.Path

	// 获得请求参数
	qs := make(map[string]string)
	for key := range c.Request.URL.Query() {
		v := c.Request.URL.Query().Get(key)
		if len(v) > 0 {
			qs[key] = v
		}
	}

	bodyRaw, _ := io.ReadAll(c.Request.Body)
	body := strings.TrimSpace(string(bodyRaw)) // 去掉换行

	// 检查签名
	return helper.Sign(signDebuggerAppID, signDebuggerAppSecret, ts, path, qs, body)
}
