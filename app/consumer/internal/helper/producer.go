package helper

import (
	"context"
	"server-api/global"
	"sync"

	"github.com/gomooth/pkg/mq"
	"github.com/gomooth/pkg/mq/kafka"
)

var (
	once     sync.Once
	producer mq.IProducer
)

func GetProducer() mq.IProducer {
	if producer != nil {
		return producer
	}

	once.Do(func() {
		producer = kafka.NewProducer(
			global.Config.Producer.Kafka.Addrs,
		)
	})
	return producer
}

// MustStartProducer 启动生产者（应在 app 启动时调用）
func MustStartProducer(ctx context.Context) error {
	p := GetProducer()
	return p.Start(ctx)
}

// MustShutdownProducer 关闭生产者（应在 app 关闭时调用）
func MustShutdownProducer(ctx context.Context) error {
	if producer != nil {
		return producer.Shutdown(ctx)
	}
	return nil
}
