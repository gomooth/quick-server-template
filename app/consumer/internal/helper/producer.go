package helper

import (
	"context"
	"server-api/global"
	"sync"

	"github.com/gomooth/pkg/mq"
	"github.com/gomooth/pkg/mq/kafka"
)

// 保留 sync.Once 保证 producer 单例；保留 shutdownOnce 保证 Shutdown 幂等。
// 不注册 global.RegisterRelease——producer 是 consumer app 的私有资源，
// 由 consumer app 的 Release/Shutdown 调用 ShutdownProducer 释放。
var (
	once         sync.Once
	shutdownOnce sync.Once
	producer     mq.IProducer
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

// StartProducer 启动生产者（应在 app 启动时调用）
func StartProducer(ctx context.Context) error {
	p := GetProducer()
	return p.Start(ctx)
}

// ShutdownProducer 关闭生产者（应在 app 关闭时调用）
func ShutdownProducer(ctx context.Context) error {
	var err error
	shutdownOnce.Do(func() {
		if producer != nil {
			err = producer.Shutdown(ctx)
		}
	})
	return err
}
