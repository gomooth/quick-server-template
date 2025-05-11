package lang

import (
	"context"
	"server-api/global"
	"server-api/service/lang"

	"github.com/gomooth/pkg/job"
)

type reloadJob struct {
}

func NewReloadJob() job.ICommandJob {
	return &reloadJob{}
}

func (s reloadJob) Run(args ...string) error {
	// 并发锁
	key := "jobTask:lang:reload"
	locker, err := global.Locker()
	if nil != err {
		return err
	}
	if err := locker.Lock(key); nil != err {
		//return xerror.Wrapf(err, "get lock failed error: [%s]", key)
		global.Log.Warningf("[%s] skip, get lock failed error: %+v", key, err)
		return nil
	}
	defer func() {
		_ = locker.UnLock(key)
	}()

	ctx := context.Background()
	return lang.Init(ctx)
}
