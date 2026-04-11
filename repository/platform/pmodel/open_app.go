package pmodel

import (
	"server-api/repository/platform/pattr"

	"gorm.io/gorm"
)

type OpenAPP struct {
	gorm.Model

	AppID        string
	AppSecret    string
	State        pattr.OpenAPPState
	PushURL      string
	PushQPSLimit uint `gorm:"column:push_qps_limit"`
}

func (OpenAPP) TableName() string {
	return "open_apps"
}
