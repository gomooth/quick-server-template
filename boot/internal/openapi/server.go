package openapi

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"server-api/global"
	"server-api/app/openapi"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomooth/pkg/framework/app"
	"github.com/gomooth/pkg/http/middleware"
)

type server struct {
	ctx        context.Context
	httpServer *http.Server
}

func New(ctx context.Context) app.IApp {
	return &server{ctx: ctx}
}

func (s *server) Start(_ context.Context) error {
	addr := strings.TrimSpace(global.Config.Server.OpenAPI.Addr)
	host := strings.TrimSpace(global.Config.Server.OpenAPI.Host)
	if len(addr) == 0 || len(host) == 0 {
		return errors.New("openapi server config invalid")
	}

	// 验证器使用 validator.v9
	//binding.Validator = new(validator.DefaultValidator)

	r := gin.New()

	// 注册全局中间件。注意顺序不要随意调整
	r.Use(gin.Recovery())
	r.Use(middleware.HttpContext())
	//r.Use(middleware.Log())
	r.Use(openapi.Middleware().IPRateLimit())
	r.Use(openapi.Middleware().RequestCache())

	// 开启 http 日志
	if global.Config.App.Log.HttpLog && global.Log != nil {
		r.Use(middleware.HttpLogger(middleware.HttpLoggerOption{
			Logger:    global.Log,
			OnlyError: global.Config.App.Log.HttpLogOnlyError,
		}))
	}

	// 非调试模式下，启用发布模式
	if global.Config.App.Log.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 注册路由
	openapi.Register(r)

	slog.Info("openapi server listening and serving HTTP", "addr", addr)
	slog.Info("openapi server starting...")
	//err := r.Run(global.Config.Server.Addr)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("openapi server start failed: %s\n", err)
		}
	}()

	return nil
}

func (s *server) Shutdown(_ context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}
