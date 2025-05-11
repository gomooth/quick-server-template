package helper

import (
	"server-api/global"
	"sync"

	"github.com/gomooth/pkg/mq/kafkaproducer"
)

var (
	once     sync.Once
	producer kafkaproducer.IProducer
)

func GetProducer() kafkaproducer.IProducer {
	if producer != nil {
		return producer
	}

	once.Do(func() {
		producer = kafkaproducer.New(
			global.Config.Producer.Kafka.Addrs,
			kafkaproducer.WithLogger(global.Log),
		)
	})
	return producer
}
