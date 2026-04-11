package consumer

import (
	"server-api/global/kafka/cg"
	"server-api/global/kafka/topic"
	"server-api/app/consumer/internal/consumer/example"

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

	return nil
}
