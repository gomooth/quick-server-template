package middleware

import (
	"fmt"
	"server-api/global"
	"time"

	"github.com/gomooth/pkg/http/middleware"

	"github.com/gin-gonic/gin"
)

func Cache() gin.HandlerFunc {
	client, err := global.RedisClient()
	if err != nil {
		panic(fmt.Sprintf("openapi cache: %v", err))
	}

	return middleware.HttpCache(
		middleware.WithHttpCacheDebug(!global.Env().IsProd()),
		middleware.WithHttpCacheLogger(global.Log),
		middleware.WithHttpCacheGlobalDuration(15*time.Minute),
		middleware.WithHttpCacheGlobalHeaderKey("X-App-Id"),
		middleware.WithHttpCacheRedisStore(client),
		middleware.WithHttpCacheGlobalSkipFields("v"),
	)
}
