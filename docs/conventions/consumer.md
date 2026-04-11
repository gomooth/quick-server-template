# 消费者约定

## NOTE 注释约定

所有注册函数中使用 `// NOTE:` 标记新功能位置，新代码加在 NOTE 之后。

## Consumer 组织

在 `app/consumer/internal/consumer/` 下按业务域创建子目录：

```
consumer/
├── example/
│   ├── redis_consumer.go
│   ├── kafka_consumer.go
│   └── httpsqs_consumer.go
└── email/
    └── redis_consumer.go
```

## 统一处理器接口

所有消费者使用 `mq.IHandler` 接口，可通过 `mq.FuncHandler` 将函数转为接口：

```go
// mq.IHandler 接口
type IHandler interface {
    Handle(ctx context.Context, msg Message) error
}

// 函数式快捷方式
var Handler mq.IHandler = mq.FuncHandler(func(ctx context.Context, msg mq.Message) error {
    // 处理消息
    return nil
})
```

`mq.Message` 包含字段：`Data []byte`、`Queue string`，以及消息队列特定方法（如 `msg.HttpsqSPosition()`）。

## Redis 消费者

```go
// consumer/example/redis_consumer.go
package example

import (
    "context"
    "log/slog"
    "github.com/gomooth/pkg/mq"
)

var RedisHandler mq.IHandler = mq.FuncHandler(func(ctx context.Context, msg mq.Message) error {
    slog.Info("redis consumer receive", "data", string(msg.Data))
    return nil
})
```

## Kafka 消费者

```go
// consumer/example/kafka_consumer.go
package example

import (
    "context"
    "log/slog"
    "github.com/gomooth/pkg/mq"
)

var KafkaConsumer mq.IHandler = mq.FuncHandler(func(ctx context.Context, msg mq.Message) error {
    slog.Debug("example kafka consumer handle, only print", "queue", msg.Queue, "msg", string(msg.Data))
    return nil
})
```

## HTTPSQS 消费者

HTTPSQS 消费者较复杂，需创建 client 和 consumer 实例：

```go
// consumer/example/httpsqs_consumer.go
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

var HttpSQSHandler mq.IHandler = mq.FuncHandler(func(ctx context.Context, msg mq.Message) error {
    pos, _ := msg.HttpsqSPosition()
    slog.Info("httpsqs consumer receive", "data", string(msg.Data), "pos", pos)
    return nil
})

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
```

## 两种注册入口

### SingleConsumerRegister — Redis/HTTPSQS

```go
// single.go
func SingleConsumerRegister(s mq.IConsumeServer) error {
    if err := s.Register("test", example.RedisHandler); err != nil {
        return err
    }

    // NOTE: 注册其它消费者

    return nil
}
```

### KafkaRegister — Kafka

```go
// kafka.go
func KafkaRegister(s mq.IConsumeServer) error {
    if err := s.Register(topic.ExampleData, example.KafkaConsumer, mq.WithGroup(cg.ExampleRecorder)); err != nil {
        return err
    }

    // NOTE: 注册其他消费者

    return nil
}
```

Kafka 注册参数顺序：`topic, handler, options...`，consumer group 通过 `mq.WithGroup()` 选项传入。

- consumer group 定义在 `global/kafka/cg/`
- topic 定义在 `global/kafka/topic/`

## Release 资源释放

每个注册入口对应一个 Release 函数：

```go
// single.go
func SingleConsumerRelease() error { return nil }

// kafka.go
func KafkaRelease() error { return nil }
```

## Kafka Producer

使用 `helper.GetProducer()` 获取单例 producer（`sync.Once` 保证）：

```go
// app/consumer/internal/helper/producer.go
pusher := helper.GetProducer()
err := pusher.Produce(ctx, topic.ExampleData, []byte("message"))
```

配置来自 `global.Config.Producer.Kafka.Addrs`。
