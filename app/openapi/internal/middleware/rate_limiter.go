package middleware

import (
	"server-api/global"

	"github.com/gin-gonic/gin"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/redis"

	libredis "github.com/redis/go-redis/v9"

	"github.com/gomooth/pkg/http/middleware"
)

const (
	ipLimit = "50-S" // 每个IP每秒50次
)

func IPRateLimit() gin.HandlerFunc {
	client := libredis.NewClient(&libredis.Options{
		Addr:     global.Config.Server.OpenAPI.Cache.Addr,
		Password: global.Config.Server.OpenAPI.Cache.Password,
		DB:       global.Config.Server.OpenAPI.Cache.DB,
	})
	store, err := redis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix: "openapi:limiter:byIp",
	})

	rate, err := limiter.NewRateFromFormatted(ipLimit)
	if err != nil {
		panic(err)
	}

	limit := limiter.New(store, rate)

	return middleware.IPRateLimit(limit)
}
