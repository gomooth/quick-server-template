package http

import (
	"server-api/app/http/internal/middleware"

	"github.com/gin-gonic/gin"
)

type customMiddleware struct{}

func Middleware() *customMiddleware {
	return &customMiddleware{}
}

func (m *customMiddleware) CORS() gin.HandlerFunc {
	return middleware.CORS()
}

func (m *customMiddleware) HTTPCache() gin.HandlerFunc {
	return middleware.HTTPCache()
}

func (m *customMiddleware) XSSFilter() gin.HandlerFunc {
	return middleware.XSSFilter()
}
