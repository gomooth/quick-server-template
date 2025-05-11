package command

import (
	"context"
	"errors"
	"time"

	"github.com/gomooth/pkg/job"

	"server-api/global"

	"github.com/save95/xerror"
)

type cmd struct {
	ctx     context.Context
	timeout time.Duration

	tasks map[string]job.ICommandJob
}

func NewCommand(ctx context.Context, timeout int) *cmd {
	return &cmd{
		ctx:     ctx,
		timeout: time.Second * time.Duration(timeout),
		tasks:   make(map[string]job.ICommandJob),
	}
}

func (c *cmd) Register(name string, cmd job.ICommandJob) {
	if _, ok := c.tasks[name]; ok {
		global.Log.Warningf("command registe skip. duplicated name: %s", name)
		return
	}

	c.tasks[name] = cmd
}

func (c *cmd) Execute(name string, args ...string) {
	task, ok := c.tasks[name]
	if !ok {
		global.Log.Warningf("command task not exist. task=%s, args=%s", name, args)
		return
	}
	if task == nil {
		global.Log.Warningf("command task is nil, skip. task=%s, args=%s", name, args)
		return
	}

	if err := task.Run(args...); nil != err {
		var xe xerror.XError
		if errors.As(err, &xe) {
			err = xe.Unwrap()
		}
		global.Log.Errorf("command failed, task=%s, args=%s, err=%+v", name, args, err)
		return
	}

	global.Log.Infof("command done, task=%s, args=%s", name, args)
}
