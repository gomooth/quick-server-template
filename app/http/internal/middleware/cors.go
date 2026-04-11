package middleware

import (
	"server-api/global"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomooth/pkg/http/middleware"
)

func CORS() gin.HandlerFunc {
	return middleware.CORS(
		middleware.WithCORSAllowOriginFunc(func(origin string) bool {
			if !global.Env().IsProd() {
				return true
			}

			// todo cors domain
			//return origin == "https://xxxx.com"
			return true
		}),
		// 分块上传请求头 https://www.npmjs.com/package/huge-uploader
		middleware.WithCORSAllowHeaders("uploader-chunk-number", "uploader-chunks-total", "uploader-file-id"),
		middleware.WithCORSHeaders("X-Custom-Key"),
		middleware.WithCORSMaxAge(24*time.Hour),
	)
}
