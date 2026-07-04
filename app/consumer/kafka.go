package consumer

import (
	"context"

	"server-api/app/consumer/internal/consumer/example"
	internalhelper "server-api/app/consumer/internal/helper"
	"server-api/global/kafka/cg"
	"server-api/global/kafka/topic"

	"github.com/gomooth/pkg/mq"
)

func KafkaRegister(s mq.IConsumeServer) error {
	if err := s.Register(topic.ExampleData, example.KafkaConsumer, mq.WithGroup(cg.ExampleRecorder)); err != nil {
		return err
	}

	// NOTE: 注册其他消费者

	return nil
}

// KafkaRelease 释放资源
func KafkaRelease() error {
	// producer 是 consumer app 的私有资源，由 app 的 Release 调用 ShutdownProducer 释放
	return internalhelper.ShutdownProducer(context.Background())
}
