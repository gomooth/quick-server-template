package lang

import (
	"context"
	"log/slog"
	"server-api/global"
	"server-api/service/lang"

	"github.com/gomooth/pkg/job"
)

type reloadJob struct {
}

func NewReloadJob() job.ICommandJob {
	return &reloadJob{}
}

func (s reloadJob) Run(_ context.Context, _ ...string) error {
	// 并发锁
	key := "jobTask:lang:reload"
	locker, err := global.Locker()
	if nil != err {
		return err
	}
	if err := locker.Lock(context.Background(), key); nil != err {
		//return xerror.Wrapf(err, "get lock failed error: [%s]", key)
		slog.Warn("skip, get lock failed", "key", key, "err", err)
		return nil
	}
	defer func() {
		_ = locker.UnLock(context.Background(), key)
	}()

	ctx := context.Background()
	return lang.Init(ctx)
}
