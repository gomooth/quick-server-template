package openapi

import (
	"server-api/app/openapi/internal/middleware"

	"github.com/gin-gonic/gin"
)

type customMiddleware struct{}

func Middleware() *customMiddleware {
	return &customMiddleware{}
}

func (m *customMiddleware) RequestCache() gin.HandlerFunc {
	return middleware.Cache()
}

func (m *customMiddleware) IPRateLimit() gin.HandlerFunc {
	return middleware.IPRateLimit()
}
