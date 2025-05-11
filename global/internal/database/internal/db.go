package internal

import "github.com/gomooth/pkg/framework/dbmanager"

type IDatabase interface {
	Platform(dbs dbmanager.IDatabaseManager) error
}
