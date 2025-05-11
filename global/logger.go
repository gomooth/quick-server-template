package global

import (
	"server-api/global/internal/logger"

	"github.com/save95/xlog"
)

var Log xlog.XLogger

func InitLogger(category string) error {
	Log = logger.New(Config.App.Log, category)
	Log.Debugf("configs: %+v", Config)
	return nil
}
