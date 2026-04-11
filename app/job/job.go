package job

import (
	"context"

	"server-api/app/job/internal/handler/example"

	"github.com/gomooth/pkg/job"
)

// CronRegister 定时任务注册
func CronRegister(r job.ICronjobRegister) {
	// 每10分钟，执行一次
	r.Register(context.Background(), "*/10 * * * *", example.NewSimpleJob())

	// NOTE: 注册其它定时任务

}

// Release 释放资源
func Release() error {

	return nil
}
