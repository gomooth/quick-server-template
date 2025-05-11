package consumer

import (
	"server-api/app/consumer/internal/consumer/example"

	"github.com/gomooth/pkg/mq/queue"
)

func SingleConsumerRegister(r queue.IRegister) {
	r.Register(example.RedisConsumer())
	//r.Register(example.HttpSQSConsumer())

	// NOTE: 注册其它消费者

}

// SingleConsumerRelease 释放资源
func SingleConsumerRelease() error {

	return nil
}
