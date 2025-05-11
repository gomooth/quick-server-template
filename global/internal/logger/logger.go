package logger

import (
	"path"

	"github.com/gomooth/pkg/framework/logger"
	"github.com/save95/xlog"
)

const defaultLogPath = "storage/logs"

func New(cnf LogConfig, category string) xlog.XLogger {
	logPath := defaultLogPath
	if len(cnf.Dir) > 0 {
		logPath = cnf.Dir
	}
	logPath = path.Join(logPath, category)

	return logger.NewLogrusLogger(
		logPath,
		logger.WithStack(xlog.DailyStack),
		logger.WithFormat(cnf.Format()),
		logger.WithStdPrint(cnf.StdPrint),
		logger.WithLevelString(cnf.Level),
	)
}
