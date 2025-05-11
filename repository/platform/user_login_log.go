package platform

import (
	"gorm.io/gorm"
)

type UserLoginLog struct {
	gorm.Model

	UserID    uint
	UserAgent string
	IP        string
	Referer   string

	UTMSource   string
	UTMMedium   string
	UTMCampaign string
	UTMTerm     string
	UTMContent  string
}
