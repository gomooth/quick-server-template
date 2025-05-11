package consumer

import (
	"context"
	"server-api/app/consumer"
	"server-api/global"
	"server-api/repository/platform"

	"github.com/gomooth/pkg/framework/app"
	"github.com/gomooth/pkg/mq/kafkaconsumer"

	"github.com/save95/xerror"
)

type kafkaServer struct {
	ctx  context.Context
	name string
	svr  app.IApp
}

func newKafka(ctx context.Context, name string) app.IApp {
	return &kafkaServer{
		ctx:  ctx,
		name: name,
	}
}

func (s *kafkaServer) Start() error {
	svr := kafkaconsumer.NewServer(
		global.Config.Consumer.Kafka.Addrs,
		kafkaconsumer.WithContext(s.ctx),
		kafkaconsumer.WithLogger(global.Log),
		kafkaconsumer.WithConsumeGroupFailedHandler(s.cgFailedHandler),
	)

	// 注册服务
	consumer.KafkaRegister(svr)

	count := svr.Count()
	if count == 0 {
		global.Log.Infof("kafka-consumer server no register consumer, skip")
		return nil
	}

	global.Log.Infof("kafka-consumer server starting, %d consumer ...", count)

	if err := svr.Start(); nil != err {
		global.Log.Errorf("kafka-consumer server start error, %s", err.Error())
		return err
	}

	// 标记在线
	flagErr := s.setRunStateFlag()

	global.Log.Infof("kafka-consumer server started. online-flag(%v)", flagErr == nil)
	s.svr = svr
	return nil
}

func (s *kafkaServer) cgFailedHandler(consumerGroup, topic string, msg []byte, err error) {
	db, derr := global.Database().Get("platform")
	if nil != derr {
		global.Log.Errorf("kafka-consumer failed saver get db failed, err=%+v", derr)
		return
	}

	record := &platform.FailedListener{
		ConsumeGroup:  consumerGroup,
		Topic:         topic,
		Msg:           string(msg),
		FailedPayload: xerror.ParsePayload(err),
		FailedReason:  xerror.FormatStackTrace(err),
	}
	if err := db.Create(record).Error; nil != err {
		global.Log.Errorf("kafka-consumer failed saver failed, err=%+v", err)
		return
	}
}

func (s *kafkaServer) getRunningFlagKey() string {
	return global.GetServerRunningFlagKey("consumer", s.name)
}

func (s *kafkaServer) setRunStateFlag() error {
	// 设置标记缓存
	key := s.getRunningFlagKey()
	redisClient, err := global.RedisClient()
	if nil != err {
		return xerror.Wrap(err, "get redis client failed")
	}
	if err := redisClient.Set(s.ctx, key, "1", 0).Err(); nil != err {
		return xerror.Wrap(err, "kafka-consumer running state flag failed")
	}
	return nil
}

func (s *kafkaServer) Shutdown() error {
	defer global.Log.Infof("kafka-consumer server stoped")

	if s.svr != nil {
		if err := s.svr.Shutdown(); err != nil {
			return xerror.Wrap(err, "kafka-consumer server shutdown failed")
		}
	}

	// 释放资源
	if err := consumer.KafkaRelease(); err != nil {
		global.Log.Errorf("kafka-consumer release failed, %+v", err)
	}

	// 清理标记缓存
	key := s.getRunningFlagKey()
	redisClient, err := global.RedisClient()
	if nil != err {
		return nil
	}
	_ = redisClient.Del(s.ctx, key).Err()

	return nil
}
