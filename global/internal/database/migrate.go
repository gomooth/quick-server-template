package database

import (
	"server-api/global/internal/config"
	"server-api/global/internal/database/internal"

	"github.com/gomooth/pkg/framework/dbmanager"
)

func Migrate(dbs dbmanager.IDatabaseManager, cnf config.ProjectConfig) error {
	if err := internal.NewMigrate().Platform(dbs); nil != err {
		return err
	}

	initer := internal.NewInit(dbs)
	if err := initer.AdminUser(cnf.App.Admin.Account, cnf.App.Admin.Password); nil != err {
		return err
	}

	return nil
}
