package consumer

import (
	"context"
	"log/slog"
	"server-api/global"
	"server-api/app/consumer"

	"github.com/gomooth/pkg/mq"
	mqredis "github.com/gomooth/pkg/mq/redis"
	goredis "github.com/redis/go-redis/v9"
)

type singleServer struct {
	svr mq.IConsumeServer
}

func newSingle() *singleServer {
	return &singleServer{}
}

func (s *singleServer) Start(ctx context.Context) error {
	rc := &global.Config.Consumer.Redis
	svr := mqredis.NewConsumer(
		rc.Addr,
		mqredis.WithConsumerRedisConfig(&goredis.Options{
			Addr:     rc.Addr,
			Password: rc.Password,
			DB:       rc.DB,
		}),
	)

	// 注册服务
	if err := consumer.SingleConsumerRegister(svr); err != nil {
		return err
	}

	count := svr.Count()
	if count == 0 {
		slog.Info("single-consumer server no register consumer, skip")
		return nil
	}

	slog.Info("single-consumer server starting", "count", count)

	if err := svr.Start(ctx); nil != err {
		slog.Error("single-consumer server start error", "err", err)
		return err
	}

	slog.Info("single-consumer server started")
	s.svr = svr
	return nil
}

func (s *singleServer) Shutdown(ctx context.Context) error {
	defer slog.Info("single-consumer server stopped")

	if s.svr != nil {
		if err := s.svr.Shutdown(ctx); err != nil {
			return err
		}
	}

	// 释放资源
	if err := consumer.SingleConsumerRelease(); err != nil {
		slog.Error("single-consumer release failed", "err", err)
	}

	return nil
}
