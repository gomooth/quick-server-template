package helper

import (
	"server-api/service/lang"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/pkg/http/restful"
)

// NewResponse 创建响应处理器，i18n 控制是否启用国际化错误提示
func NewResponse(ctx *gin.Context, i18n bool, opts ...restful.ResponseOption) restful.IResponse {
	if i18n {
		opts = append([]restful.ResponseOption{restful.WithResponseErrorMsgHandler(lang.Handler())}, opts...)
	}

	return restful.NewResponse(ctx, opts...)
}
