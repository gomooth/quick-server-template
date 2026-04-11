package example

import (
	"context"
	"log/slog"
	"server-api/global"
	"time"

	httpsqsdirect "github.com/gomooth/httpsqs"
	"github.com/gomooth/pkg/mq"
	"github.com/gomooth/pkg/mq/httpsqs"
)

// HttpSQSHandler HTTPSQS 消费者处理器
var HttpSQSHandler mq.IHandler = mq.FuncHandler(func(ctx context.Context, msg mq.Message) error {
	pos, _ := msg.HttpsqSPosition()
	slog.Info("httpsqs consumer receive", "data", string(msg.Data), "pos", pos)
	return nil
})

// NewHTTPSQSConsumer 创建 HTTPSQS 消费者
func NewHTTPSQSConsumer() mq.IConsumeServer {
	client := httpsqsdirect.NewClient(&httpsqsdirect.Config{
		Addr:     global.Config.Consumer.HTTPSQS.Addr,
		Password: global.Config.Consumer.HTTPSQS.Password,
		Timeout:  time.Duration(global.Config.Consumer.HTTPSQS.Timeout) * time.Second,
	})

	return httpsqs.NewConsumer(
		httpsqs.WithHTTPSQSClient(client),
		httpsqs.WithMaxRetry(3),
		httpsqs.WithConsumer("test", HttpSQSHandler),
	)
}
