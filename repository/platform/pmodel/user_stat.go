package pmodel

import (
	"time"

	"server-api/repository/platform/pattr"

	"gorm.io/gorm"
)

type UserStat struct {
	gorm.Model

	UserID         uint
	InviterID      uint
	FromGenre      uint
	FromPlatformID pattr.UserFromPlatform
	FromChannelID  uint
	UTMSource      string
	UTMMedium      string
	UTMCampaign    string
	UTMTerm        string
	UTMContent     string
	LastVisitAt    *time.Time
}
