package platform

import (
	"time"

	"gorm.io/gorm"
)

type UserLoginLog struct {
	ID uint `gorm:"type:INT(11) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`

	UserID      uint   `gorm:"type:INT(11) UNSIGNED;not null;index:idx_user"`
	UserAgent   string `gorm:"not null;size:512"`
	IP          string `gorm:"not null;size:32"`
	Referer     string `gorm:"not null;size:256"`
	UTMSource   string `gorm:"column:utm_source;size:256"`
	UTMMedium   string `gorm:"column:utm_medium;size:256"`
	UTMCampaign string `gorm:"column:utm_campaign;size:256"`
	UTMTerm     string `gorm:"column:utm_term;size:256"`
	UTMContent  string `gorm:"column:utm_content;size:256"`

	CreatedAt time.Time      `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time      `gorm:"not null;default:current_timestamp on update current_timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_deleted_at"`
}
