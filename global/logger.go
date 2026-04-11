package global

import (
	"log/slog"
	"server-api/global/internal/logger"
)

var Log *slog.Logger

func InitLogger(category string) error {
	Log = logger.New(Config.App.Log, category)
	slog.SetDefault(Log)
	slog.Debug("configs", "config", Config)
	return nil
}
