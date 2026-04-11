package middleware

import (
	"server-api/global"
	"server-api/app/http/internal/helper"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomooth/pkg/http/jwt"
	"github.com/gomooth/pkg/http/middleware"
	"github.com/redis/go-redis/v9"
)

func HTTPCache() gin.HandlerFunc {
	jwtOpt := helper.JWTOption(false)
	return middleware.HttpCache(
		middleware.WithHttpCacheLogger(global.Log),
		middleware.WithHttpCacheUserIDFunc(func(c *gin.Context) (uint, error) {
			user, err := jwt.ParseJWTUser(c, jwtOpt)
			if err != nil || user == nil {
				return 0, err
			}
			return user.ID, nil
		}),
		middleware.WithHttpCacheGlobalDuration(5*time.Minute),
		middleware.WithHttpCacheRedisStore(redis.NewClient(&redis.Options{
			Addr:     global.Config.Server.HTTP.Cache.Addr,
			Password: global.Config.Server.HTTP.Cache.Password,
			DB:       global.Config.Server.HTTP.Cache.DB,
		})),
		middleware.WithHttpCacheGlobalSkipFields("v"),
		middleware.WithHttpCacheRouteSkipFiledPolicy("/user/", false),
	)
}