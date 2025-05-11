package example

import (
	"fmt"
	"server-api/app/job/internal/helper"
	"server-api/global"

	"github.com/gomooth/pkg/job"
)

type simpleJob struct {
}

func NewSimpleJob() job.ICommandJob {
	return &simpleJob{}
}

func (s simpleJob) Run(args ...string) error {
	global.Log.Debugf("example simple job, only print. args=%#v", args)

	params := helper.NewCMDArgs(args...)
	version := params.Get("ver", "version")
	isTest := params.GetBool("test")
	fmt.Printf("example simple job args: version=%s, isTest=%v\n", version, isTest)

	return nil
}
