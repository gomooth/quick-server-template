package boot

import (
	"context"
	"server-api/boot/internal/command"

	jobapp "server-api/app/job"
	"server-api/global"

	"github.com/pkg/errors"
)

func Command(cnf Param) error {
	if err := initialize(cnf); nil != err {
		return errors.Wrap(err, "initialize failed")
	}

	conf := cnf.CMDParam
	ctx := context.Background()
	cmd := command.NewCommand(ctx, conf.Timeout)

	// 注册所有命令
	jobapp.CMDRegister(cmd)

	// 执行命令
	cmd.Execute(conf.Name, conf.Args...)

	global.Log.Infof("[%s] command done", conf.Name)
	return nil
}
