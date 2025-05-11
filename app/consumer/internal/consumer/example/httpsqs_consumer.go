package example

import (
	"context"
	"server-api/global"
	"time"

	"github.com/gomooth/httpsqs"
	"github.com/gomooth/pkg/mq/httpsqsconsumer"
	"github.com/gomooth/pkg/mq/queue"
)

func HttpSQSConsumer() queue.IConsumer {
	return httpsqsconsumer.New(
		httpsqsconsumer.WithLogger(global.Log),
		httpsqsconsumer.WithMaxRetry(3),
		httpsqsconsumer.WithHandler(&httpSQSConsumerHandler{}),
	)
}

type httpSQSConsumerHandler struct{}

func (h *httpSQSConsumerHandler) QueueName() string {
	return "test"
}

func (h *httpSQSConsumerHandler) GetClient() (httpsqs.IClient, error) {
	return httpsqs.NewClient(&httpsqs.Config{
		Addr:     global.Config.Consumer.HTTPSQS.Addr,
		Password: global.Config.Consumer.HTTPSQS.Password,
		Timeout:  time.Duration(global.Config.Consumer.HTTPSQS.Timeout) * time.Second,
	}), nil
}

func (h *httpSQSConsumerHandler) OnBefore(ctx context.Context) error {
	return nil
}

func (h *httpSQSConsumerHandler) Handle(ctx context.Context, data string, pos int64) error {
	global.Log.Infof("[consumer] httpsqs consumer receive: %s, pos: %d", data, pos)
	return nil
}

func (h *httpSQSConsumerHandler) OnFailed(ctx context.Context, data string, err error) {
}
