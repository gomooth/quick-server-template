package logger

import (
	"strings"

	"github.com/save95/xlog"
)

type LogConfig struct {
	Dir              string `toml:"dir"`                 // 存储目录路径
	Level            string `toml:"level"`               // 日志级别
	FormatText       string `toml:"format"`              // 日志格式：text-文本(默认)，json-JSON(普遍应用于云服务器)
	StdPrint         bool   `toml:"std_print"`           // 是否打印到控制台
	HttpLog          bool   `toml:"http_log"`            // 是否打印 http 日志
	HttpLogOnlyError bool   `toml:"http_log_only_error"` // 是否仅打印 http 错误日志
}

func (m LogConfig) Format() xlog.LogFormat {
	switch strings.ToLower(m.FormatText) {
	case "json":
		return xlog.LogFormatJson
	case "text":
		return xlog.LogFormatText
	default:
		return xlog.LogFormatText
	}
}
