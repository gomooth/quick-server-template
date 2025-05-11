package middleware

import (
	"server-api/app/http/internal/helper"
	"server-api/global"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomooth/pkg/http/middleware"
	"github.com/redis/go-redis/v9"
)

func HTTPCache() gin.HandlerFunc {
	return middleware.HttpCache(
		//middleware.WithHttpCacheDebug(!global.Env().IsProd()),
		middleware.WithHttpCacheLogger(global.Log),
		middleware.WithHttpCacheJWTOption(helper.JWTOption(false)),
		middleware.WithHttpCacheGlobalDuration(5*time.Minute),
		middleware.WithHttpCacheRedisStore(redis.NewClient(&redis.Options{
			Addr:     global.Config.Server.HTTP.Cache.Addr,
			Password: global.Config.Server.HTTP.Cache.Password,
			DB:       global.Config.Server.HTTP.Cache.DB,
		})),
		middleware.WithHttpCacheGlobalSkipFields("v"),
		middleware.WithHttpCacheRouteSkipFiledPolicy("/user/", true),
	)
}
