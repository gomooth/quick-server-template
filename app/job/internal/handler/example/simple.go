package example

import (
	"context"
	"fmt"
	"log/slog"
	"server-api/app/job/internal/helper"

	"github.com/gomooth/pkg/job"
)

type simpleJob struct {
}

func NewSimpleJob() job.ICommandJob {
	return &simpleJob{}
}

func (s simpleJob) Run(_ context.Context, args ...string) error {
	slog.Debug("example simple job, only print", "args", args)

	params := helper.NewCMDArgs(args...)
	version := params.Get("ver", "version")
	isTest := params.GetBool("test")
	fmt.Printf("example simple job args: version=%s, isTest=%v\n", version, isTest)

	return nil
}
