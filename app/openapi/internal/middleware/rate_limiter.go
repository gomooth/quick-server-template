package middleware

import (
	"fmt"
	"server-api/global"

	"github.com/gin-gonic/gin"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/redis"

	"github.com/gomooth/pkg/http/middleware"
)

const (
	ipLimit = "50-S" // 每个IP每秒50次
)

func IPRateLimit() gin.HandlerFunc {
	client, err := global.RedisClient()
	if err != nil {
		panic(fmt.Sprintf("openapi rate limiter: %v", err))
	}
	store, err := redis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix: "openapi:limiter:byIp",
	})
	if err != nil {
		panic(fmt.Sprintf("openapi rate limiter store: %v", err))
	}

	rate, err := limiter.NewRateFromFormatted(ipLimit)
	if err != nil {
		panic(err)
	}

	limit := limiter.New(store, rate)

	return middleware.IPRateLimit(limit)
}
