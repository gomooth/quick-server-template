package example

import (
	"context"
	"log/slog"

	"github.com/gomooth/pkg/mq"
)

// KafkaConsumer Kafka 消费者处理器
var KafkaConsumer mq.IHandler = mq.FuncHandler(func(ctx context.Context, msg mq.Message) error {
	slog.Debug("example kafka consumer handle, only print", "queue", msg.Queue, "msg", string(msg.Data))
	return nil
})
