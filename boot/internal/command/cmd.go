package command

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gomooth/pkg/job"

	"github.com/gomooth/xerror"
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

func (c *cmd) Register(_ context.Context, name string, cmd job.ICommandJob) {
	if _, ok := c.tasks[name]; ok {
		slog.Warn("command registe skip", "name", name)
		return
	}

	c.tasks[name] = cmd
}

func (c *cmd) Execute(name string, args ...string) {
	task, ok := c.tasks[name]
	if !ok {
		slog.Warn("command task not exist", "task", name, "args", args)
		return
	}
	if task == nil {
		slog.Warn("command task is nil, skip", "task", name, "args", args)
		return
	}

	if err := task.Run(c.ctx, args...); nil != err {
		var xe xerror.XError
		if errors.As(err, &xe) {
			err = xe.Unwrap()
		}
		slog.Error("command failed", "task", name, "args", args, "err", err)
		return
	}

	slog.Info("command done", "task", name, "args", args)
}
