package database

import (
	"fmt"

	"github.com/gomooth/pkg/framework/dbmanager"
	"github.com/save95/xerror"

	"gorm.io/gorm"
)

var dbs = make(map[string]*gorm.DB)

type databases struct {
}

func Database() dbmanager.IDatabaseManager {
	return &databases{}
}

func (db databases) Register(name string, dbc *gorm.DB) error {
	if len(name) == 0 || dbc == nil {
		return xerror.New("database register params error")
	}

	if _, ok := dbs[name]; ok {
		return xerror.New(fmt.Sprintf("%s database duplicate registration", name))
	}

	dbs[name] = dbc
	return nil
}

func (db databases) Get(name string) (*gorm.DB, error) {
	c, ok := dbs[name]
	if !ok {
		return nil, xerror.New(fmt.Sprintf("%s database not registered", name))
	}

	return c, nil
}
