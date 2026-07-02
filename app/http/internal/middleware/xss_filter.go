package middleware

import (
	"server-api/global"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/pkg/http/middleware"
	"github.com/gomooth/pkg/http/xss"
)

func XSSFilter() gin.HandlerFunc {
	return middleware.XSSFilter(
		middleware.WithXSSDebug(!global.Env().IsProd()),
		middleware.WithXSSGlobalPolicy(xss.PolicyStrict),
		middleware.WithXSSGlobalSkipFields("password"),
		middleware.WithXSSRoutePolicy("admin", xss.PolicyUGC),
		middleware.WithXSSRoutePolicy("/callback/", xss.PolicyNone),
		middleware.WithXSSRoutePolicy("/endpoint", xss.PolicyNone),
		middleware.WithXSSRoutePolicy("/ping", xss.PolicyNone),
		middleware.WithXSSRoutePolicy("/healthz", xss.PolicyNone),
		middleware.WithXSSRoutePolicy("/readyz", xss.PolicyNone),
	)
}
