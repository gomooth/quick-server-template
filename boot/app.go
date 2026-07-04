package boot

import (
	"context"
	"log/slog"
	"server-api/boot/internal/consumer"
	"server-api/boot/internal/cronjob"
	"server-api/boot/internal/http"
	"server-api/boot/internal/openapi"
	"server-api/boot/internal/watcher"
	"server-api/global"
	"server-api/service/lang"
	"time"

	"github.com/gomooth/pkg/framework/app"

	"github.com/gomooth/xerror"

	"github.com/fsnotify/fsnotify"

	"github.com/pkg/errors"
)

func initialize(cnf Param) error {
	// 加载配置
	if err := global.ParseConfig(cnf.ConfigFilename); nil != err {
		return errors.Wrap(err, "parser config file failed")
	}

	// 初始化日志
	if err := global.InitLogger(cnf.LogCategory()); err != nil {
		return errors.Wrap(err, "init logger failed")
	}

	// 初始化db
	if err := global.InitDataBase(); err != nil {
		return errors.Wrap(err, "init db connect failed")
	}

	return nil
}

func Boot(cnf Param) error {
	if err := initialize(cnf); nil != err {
		return errors.Wrap(err, "initialize failed")
	}

	ctx := context.Background()

	// 初始化数据
	if err := global.MigrateData(); err != nil {
		return errors.Wrap(err, "init db connect failed")
	}

	// 初始化超级管理员
	if err := seedAdminUser(global.Config.App.Admin.Account, global.Config.App.Admin.Password); err != nil {
		return errors.Wrap(err, "seed admin user failed")
	}

	// 初始化语言包
	if err := lang.Init(ctx); nil != err {
		return xerror.Wrap(err, "lang init failed")
	}

	// 注册 apps
	apps := app.NewManager(
		app.WithLogger(global.Log),
	)

	// 注册 配置文件监听器
	if global.Config.App.WatchConfigEnabled {
		slog.Debug("watch config file charge enabled")
		localname, _ := global.GetConfigFilename(cnf.ConfigFilename)
		apps.Register(watcher.NewFileServer(ctx, localname, func(ev fsnotify.Event) error {
			// 配置文件被修改，则更新全局配置
			if ev.Op == fsnotify.Write {
				return global.ParseConfig(cnf.ConfigFilename)
			}
			return nil
		}))
	}

	// 注册应用服务
	for _, server := range cnf.RegisterServers {
		switch server {
		case InitServerTypeWeb:
			apps.Register(http.New(ctx))
		case InitServerTypeOpenAPI:
			apps.Register(openapi.New(ctx))
		case InitServerTypeCronjob:
			apps.Register(cronjob.New(ctx))
		case InitServerTypeConsumer:
			apps.Register(consumer.New(ctx))
		}
	}

	// 阻塞直至收到中断信号或启动失败返回；返回值含启动/关闭过程中的错误
	runErr := apps.Run(ctx)

	// 无论 Run 成败，都释放已注册的基础设施资源（Redis/Cache/DB/Producer 等），LIFO 顺序
	releaseCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := global.Release(releaseCtx); err != nil {
		slog.Error("global release failed", "err", err)
	}
	return runErr
}
