package job

import (
	"context"

	"server-api/app/job/internal/handler/example"
	"server-api/app/job/internal/handler/lang"

	"github.com/gomooth/pkg/job"
)

func CMDRegister(r job.ICommandRegister) {
	r.Register(context.Background(), "example-simple", example.NewSimpleJob())

	// NOTE: 注册其它命令

	r.Register(context.Background(), "lang:reload", lang.NewReloadJob())

}
