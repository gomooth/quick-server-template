package platform

import (
	"time"

	"gorm.io/gorm"
)

type UserStat struct {
	ID uint `gorm:"type:INT(11) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`

	UserID         uint       `gorm:"type:INT(11) UNSIGNED;not null;uniqueIndex:uk_user"`
	InviterID      uint       `gorm:"type:INT(11) UNSIGNED;not null;comment:邀请者ID"`
	FromGenre      uint       `gorm:"type:INT(11) UNSIGNED;not null;comment:来源方式：0-自然流量"`
	FromPlatformID uint       `gorm:"type:INT(11) UNSIGNED;not null;comment:来源平台ID"`
	FromChannelID  uint       `gorm:"type:INT(11) UNSIGNED;not null;comment:来源渠道ID"`
	UTMSource      string     `gorm:"column:utm_source;size:256;comment:注册时广告投放来源"`
	UTMMedium      string     `gorm:"column:utm_medium;size:256;comment:注册时广告投放媒介"`
	UTMCampaign    string     `gorm:"column:utm_campaign;size:256;comment:注册时广告投放名称"`
	UTMTerm        string     `gorm:"column:utm_term;size:256;comment:注册时广告投放字词关键字"`
	UTMContent     string     `gorm:"column:utm_content;size:256;comment:注册时广告投放内容"`
	LastVisitAt    *time.Time `gorm:"comment:最后访问时间"`

	CreatedAt time.Time      `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time      `gorm:"not null;default:current_timestamp on update current_timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_deleted_at"`
}
