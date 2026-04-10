package global

import (
	"server-api/global/internal/logger"

	"github.com/save95/xlog"
)

var Log xlog.XLogger

func InitLogger(category string) error {
	Log = logger.New(Config.App.Log, category)
	Log.Debugf("configs: %+v", Config)

	// 初始化 slog 默认 logger，使业务侧可以直接使用 slog 包
	logger.ActiveSlog(Log)

	return nil
}
