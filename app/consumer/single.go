package consumer

import (
	"server-api/app/consumer/internal/consumer/example"

	"github.com/gomooth/pkg/mq"
)

func SingleConsumerRegister(s mq.IConsumeServer) error {
	if err := s.Register("test", example.RedisHandler); err != nil {
		return err
	}

	// NOTE: 注册其它消费者

	return nil
}

// SingleConsumerRelease 释放资源
func SingleConsumerRelease() error {

	return nil
}
