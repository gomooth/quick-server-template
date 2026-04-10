package system

import (
	"io"
	"server-api/app/openapi/internal/helper"
	"server-api/repository/platform/dao"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/pkg/http/restful"

	"github.com/save95/xerror"
)

type Controller struct{}

func (c Controller) CurrentTime(ctx *gin.Context) {
	rru := restful.NewResponse(ctx)

	now := time.Now()
	rru.Retrieve(map[string]interface{}{
		"zone":      time.Now().In(time.Local).Location().String(),
		"time":      now.Format("2006-01-02 15:04:05.999"),
		"timestamp": now.UnixMilli(),
	})
}

func (c Controller) Sign(ctx *gin.Context) {
	rru := restful.NewResponse(ctx)

	h, err := helper.MustParseHeader(ctx)
	if err != nil {
		rru.WithError(err)
		return
	}

	app, err := dao.NewOpenAPP().FirstByAppID(ctx, h.AppID)
	if nil != err {
		rru.WithError(xerror.New("获取应用信息失败"))
		return
	}

	path := ctx.Request.URL.Path

	// 获得请求参数
	qs := make(map[string]string)
	for key := range ctx.Request.URL.Query() {
		v := ctx.Request.URL.Query().Get(key)
		if len(v) > 0 {
			qs[key] = v
		}
	}

	bodyRaw, _ := io.ReadAll(ctx.Request.Body)
	body := strings.TrimSpace(string(bodyRaw)) // 去掉换行

	result := helper.Sign(h.AppID, app.AppSecret, h.Timestamp, path, qs, body)
	result.Input = ctx.GetHeader("X-Sign")
	result.Success = result.Input == result.Signature
	rru.Retrieve(result)
}
