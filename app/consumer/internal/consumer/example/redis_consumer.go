package example

import (
	"context"
	"log/slog"

	"github.com/gomooth/pkg/mq"
)

// RedisHandler Redis 消费者处理器
var RedisHandler mq.IHandler = mq.FuncHandler(func(ctx context.Context, msg mq.Message) error {
	slog.Info("redis consumer receive", "data", string(msg.Data))
	return nil
})
