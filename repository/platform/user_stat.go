package platform

import (
	"time"

	"server-api/repository/types/platformtypes"

	"gorm.io/gorm"
)

type UserStat struct {
	gorm.Model

	UserID         uint
	InviterID      uint
	FromGenre      uint
	FromPlatformID platformtypes.UserFromPlatform
	FromChannelID  uint
	UTMSource      string
	UTMMedium      string
	UTMCampaign    string
	UTMTerm        string
	UTMContent     string
	LastVisitAt    *time.Time
}
