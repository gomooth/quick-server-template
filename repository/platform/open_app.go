package platform

import (
	"server-api/repository/types/platformtypes"

	"gorm.io/gorm"
)

type OpenAPP struct {
	gorm.Model

	AppID        string
	AppSecret    string
	State        platformtypes.OpenAPPState
	PushURL      string
	PushQPSLimit uint `gorm:"column:push_qps_limit"`
}

func (OpenAPP) TableName() string {
	return "open_apps"
}
