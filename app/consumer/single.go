package consumer

import (
	"context"

	"server-api/app/consumer/internal/consumer/example"
	internalhelper "server-api/app/consumer/internal/helper"

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
	// producer 是 consumer app 的私有资源，由 app 的 Release 调用 ShutdownProducer 释放
	return internalhelper.ShutdownProducer(context.Background())
}
