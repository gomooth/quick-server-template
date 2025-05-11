package consumer

import (
	"context"
	"server-api/app/consumer"
	"server-api/global"

	"github.com/gomooth/pkg/framework/app"
	"github.com/gomooth/pkg/mq/queue"
)

type singleServer struct {
	ctx context.Context
	svr app.IApp
}

func newSingle(ctx context.Context) app.IApp {
	return &singleServer{
		ctx: ctx,
	}
}

func (s *singleServer) Start() error {
	svr := queue.NewServer(s.ctx)

	// 注册服务
	consumer.SingleConsumerRegister(svr)

	count := svr.Count()
	if count == 0 {
		global.Log.Infof("single-consumer server no register consumer, skip")
		return nil
	}

	global.Log.Infof("single-consumer server starting, %d consumer ...", count)

	if err := svr.Start(); nil != err {
		global.Log.Errorf("single-consumer server start error, %s", err.Error())
		return err
	}

	global.Log.Info("single-consumer server started")
	s.svr = svr
	return nil
}

func (s *singleServer) Shutdown() error {
	defer global.Log.Infof("single-consumer server stoped")

	if s.svr != nil {
		if err := s.svr.Shutdown(); nil != err {
			return err
		}
	}

	// 释放资源
	if err := consumer.SingleConsumerRelease(); err != nil {
		global.Log.Errorf("single-consumer release failed, %+v", err)
	}

	return nil
}
