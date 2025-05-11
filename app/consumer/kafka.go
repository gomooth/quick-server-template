package consumer

import (
	"server-api/app/consumer/internal/consumer/example"
	"server-api/global/kafka/cg"
	"server-api/global/kafka/topic"

	"github.com/gomooth/pkg/mq/kafkaconsumer"
)

func KafkaRegister(s kafkaconsumer.IRegister) {
	s.Register(cg.ExampleRecorder, example.KafkaConsumer, topic.ExampleData)

	// NOTE: 注册其他消费者

}

// KafkaRelease 释放资源
func KafkaRelease() error {

	return nil
}
