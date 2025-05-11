package example

import (
	"server-api/global"

	"github.com/gomooth/pkg/mq/queue"
	"github.com/gomooth/pkg/mq/redisconsumer"
)

var cnf = &queue.RedisQueueConfig{
	Addr:     global.Config.Consumer.Redis.Addr,
	Password: global.Config.Consumer.Redis.Password,
}

func RedisConsumer() queue.IConsumer {
	return redisconsumer.New(
		redisconsumer.WithLogger(global.Log),
		redisconsumer.WithHandler(cnf, "", func(val string) error {
			global.Log.Infof("[consumer] redis consumer receive: %s", val)
			return nil
		}))
}
