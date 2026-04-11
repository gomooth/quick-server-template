package logger

import (
	"log/slog"
	"path"

	pkglogger "github.com/gomooth/pkg/framework/logger"
)

const defaultLogPath = "storage/logs"

func New(cnf LogConfig, category string) *slog.Logger {
	logPath := defaultLogPath
	if len(cnf.Dir) > 0 {
		logPath = cnf.Dir
	}
	logPath = path.Join(logPath, category)

	return pkglogger.NewFileLogger(
		logPath,
		pkglogger.WithLevelString(cnf.Level),
		pkglogger.WithStdPrint(cnf.StdPrint),
	)
}
