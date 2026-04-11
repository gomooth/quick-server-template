package consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"server-api/global"
	"server-api/app/consumer"
	"server-api/repository/platform/pmodel"

	"github.com/gomooth/pkg/mq"
	"github.com/gomooth/pkg/mq/kafka"

	"github.com/gomooth/xerror"
)

type kafkaServer struct {
	name string
	svr  mq.IConsumeServer
}

func newKafka(name string) *kafkaServer {
	return &kafkaServer{
		name: name,
	}
}

func (s *kafkaServer) Start(ctx context.Context) error {
	svr := kafka.NewConsumer(
		global.Config.Consumer.Kafka.Addrs,
		kafka.WithFailedHandler(s.failedHandler),
	)

	// 注册服务
	if err := consumer.KafkaRegister(svr); err != nil {
		return xerror.Wrap(err, "kafka consumer register failed")
	}

	count := svr.Count()
	if count == 0 {
		slog.Info("kafka-consumer server no register consumer, skip")
		return nil
	}

	slog.Info("kafka-consumer server starting", "count", count)

	if err := svr.Start(ctx); nil != err {
		slog.Error("kafka-consumer server start error", "err", err)
		return err
	}

	// 标记在线
	flagErr := s.setRunStateFlag(ctx)

	slog.Info("kafka-consumer server started", "online", flagErr == nil)
	s.svr = svr
	return nil
}

func (s *kafkaServer) failedHandler(ctx context.Context, msg mq.Message, err error) {
	db, derr := global.Database().Get("platform")
	if nil != derr {
		slog.Error("kafka-consumer failed saver get db failed", "err", derr)
		return
	}

	group, _ := msg.KafkaGroup()
	record := &pmodel.FailedListener{
		ConsumeGroup:  group,
		Topic:         msg.Queue,
		Msg:           string(msg.Data),
		FailedPayload: func() string {
			payload := xerror.ParsePayload(err)
			bs, _ := json.Marshal(payload)
			return string(bs)
		}(),
		FailedReason: xerror.StackTrace(err),
	}
	if err := db.Create(record).Error; nil != err {
		slog.Error("kafka-consumer failed saver failed", "err", err)
		return
	}
}

func (s *kafkaServer) getRunningFlagKey() string {
	return global.GetServerRunningFlagKey("consumer", s.name)
}

func (s *kafkaServer) setRunStateFlag(ctx context.Context) error {
	// 设置标记缓存
	key := s.getRunningFlagKey()
	redisClient, err := global.RedisClient()
	if nil != err {
		return xerror.Wrap(err, "get redis client failed")
	}
	if err := redisClient.Set(ctx, key, "1", 0).Err(); nil != err {
		return xerror.Wrap(err, "kafka-consumer running state flag failed")
	}
	return nil
}

func (s *kafkaServer) Shutdown(ctx context.Context) error {
	defer slog.Info("kafka-consumer server stopped")

	if s.svr != nil {
		if err := s.svr.Shutdown(ctx); err != nil {
			return xerror.Wrap(err, "kafka-consumer server shutdown failed")
		}
	}

	// 释放资源
	if err := consumer.KafkaRelease(); err != nil {
		slog.Error("kafka-consumer release failed", "err", err)
	}

	// 清理标记缓存
	key := s.getRunningFlagKey()
	redisClient, err := global.RedisClient()
	if nil != err {
		return nil
	}
	_ = redisClient.Del(ctx, key).Err()

	return nil
}
