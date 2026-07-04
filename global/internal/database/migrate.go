package database

import (
	"server-api/global/internal/database/internal"

	"github.com/gomooth/pkg/framework/dbmanager"
)

func Migrate(dbs dbmanager.IDatabaseManager) error {
	if err := internal.NewMigrate().Platform(dbs); nil != err {
		return err
	}

	return nil
}
