package middleware

import (
	"server-api/global"
	"time"

	"github.com/gomooth/pkg/http/middleware"

	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
)

func Cache() gin.HandlerFunc {
	return middleware.HttpCache(
		middleware.WithHttpCacheDebug(!global.Env().IsProd()),
		middleware.WithHttpCacheLogger(global.Log),
		middleware.WithHttpCacheGlobalDuration(15*time.Minute),
		middleware.WithHttpCacheGlobalHeaderKey("X-App-Id"),
		middleware.WithHttpCacheRedisStore(redis.NewClient(&redis.Options{
			Addr:     global.Config.Server.OpenAPI.Cache.Addr,
			Password: global.Config.Server.OpenAPI.Cache.Password,
			DB:       global.Config.Server.OpenAPI.Cache.DB,
		})),
		middleware.WithHttpCacheGlobalSkipFields("v"),
	)
}