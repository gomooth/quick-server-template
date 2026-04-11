package database

import (
	"fmt"

	"github.com/gomooth/pkg/framework/dbmanager"
	"github.com/gomooth/xerror"

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

func (db databases) Unregister(name string) error {
	if _, ok := dbs[name]; !ok {
		return xerror.New(fmt.Sprintf("%s database not registered", name))
	}
	delete(dbs, name)
	return nil
}

func (db databases) CloseAll() error {
	var errs []error
	for name, dbc := range dbs {
		if sqlDB, err := dbc.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errs = append(errs, xerror.Wrap(err, fmt.Sprintf("close %s database failed", name)))
			}
		}
		delete(dbs, name)
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
