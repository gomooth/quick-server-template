package cronjob

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"server-api/global"
	jobapp "server-api/app/job"
	"server-api/repository/platform/pmodel"
	"strings"

	"github.com/gomooth/xerror"

	"github.com/gomooth/pkg/framework/app"
	"github.com/gomooth/pkg/job"

	"github.com/robfig/cron/v3"
)

type server struct {
	ctx context.Context
	c   *cron.Cron
}

func New(ctx context.Context) app.IApp {
	return &server{
		ctx: ctx,
		c:   cron.New(),
	}
}

func (s server) Start(_ context.Context) error {
	slog.Info("cronjob server starting...")

	// 注册定时脚本
	jobapp.CronRegister(s)

	s.c.Start()
	slog.Info("cronjob server started")
	return nil
}

func (s server) Register(_ context.Context, spec string, cmd job.ICommandJob) {
	wrapper := job.NewCronJobWrapper(
		job.WrapWithLogger(global.Log),
		job.WrapWithFailedSaver(s.failedSaver),
	)

	eid, err := s.c.AddJob(spec, wrapper.FromCommandJob(s.ctx, cmd))
	if nil != err {
		slog.Error("cronjob register failed", "err", err)
		return
	}

	name := strings.Trim(fmt.Sprintf("%T", cmd), "*")
	slog.Info("cronjob register success", "name", name, "entryID", int(eid))
	return
}

func (s server) failedSaver(jobName string, in []string, err error) {
	db, derr := global.Database().Get("platform")
	if nil != derr {
		slog.Error("cronjob failed saver get db failed", "err", derr)
		return
	}

	argsBytes, _ := json.Marshal(in)

	payload := xerror.ParsePayload(err)
	payloadBytes, _ := json.Marshal(payload)

	record := &pmodel.FailedJob{
		JobName:     jobName,
		JobArgs:     string(argsBytes),
		Payload:     string(payloadBytes),
		Errors:      xerror.StackTrace(err),
		Handled:     false,
		HandledAt:   nil,
		Compensated: false,
	}
	if err := db.Create(record).Error; nil != err {
		slog.Error("cronjob failed saver failed", "err", err)
		return
	}
}

func (s server) Shutdown(_ context.Context) error {
	defer slog.Info("cronjob server stop")

	if s.c != nil {
		s.c.Stop()
	}

	// 释放资源
	if err := jobapp.Release(); err != nil {
		slog.Error("cronjob release failed", "err", err)
	}

	return nil
}
