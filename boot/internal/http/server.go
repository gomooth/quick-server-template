package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	httpapp "server-api/app/http"
	"server-api/global"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/pkg/framework/app"
	"github.com/gomooth/pkg/http/middleware"

	"github.com/save95/xerror"
	"github.com/save95/xlog"
)

type server struct {
	ctx        context.Context
	httpServer *http.Server
}

func New(ctx context.Context) app.IApp {
	return &server{ctx: ctx}
}

func (s *server) Start() error {
	addr := strings.TrimSpace(global.Config.Server.HTTP.Addr)
	host := strings.TrimSpace(global.Config.Server.HTTP.Host)
	if len(addr) == 0 || len(host) == 0 {
		return errors.New("http server config invalid")
	}

	// 验证器使用 validator.v9
	//binding.Validator = new(validator.DefaultValidator)

	r := gin.New()

	// 注册全局中间件。注意顺序不要随意调整
	r.Use(gin.Recovery())
	r.Use(httpapp.Middleware().CORS())

	// 开启 http 缓存
	if global.Config.Server.HTTP.Cache.Enabled {
		r.Use(httpapp.Middleware().HTTPCache())
	}

	r.Use(httpapp.Middleware().XSSFilter())
	r.Use(middleware.HttpContext())
	//r.Use(middleware.RESTFul(global.ApiVersionLatest))
	//r.Use(middleware.Log())
	//r.Use(middleware.ParserSession())

	// 开启 http 日志
	if global.Config.App.Log.HttpLog && global.Log != nil {
		r.Use(middleware.HttpLogger(middleware.HttpLoggerOption{
			Logger:    global.Log,
			OnlyError: global.Config.App.Log.HttpLogOnlyError,
		}))
	}

	// 非调试模式下，启用发布模式
	if global.Env().IsProd() && xlog.ParseLevel(global.Config.App.Log.Level) != xlog.DebugLevel {
		gin.SetMode(gin.ReleaseMode)
	}

	// 注册路由，路由统一安置在 app/http/route 目录，由 main 引导
	httpapp.RouteRegister(r)

	global.Log.Infof("http server listening and serving HTTP on %s", addr)
	global.Log.Info("http server starting...")

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server start failed: %s\n", err)
		}
	}()

	return nil
}

func (s *server) Shutdown() error {
	defer global.Log.Infof("http server stop")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); nil != err {
		return xerror.Wrap(err, "http server shutdown failed")
	}

	// 释放资源
	if err := httpapp.Release(); err != nil {
		global.Log.Errorf("http server release failed, %+v", err)
	}

	return nil
}
