package platform

import (
	"time"

	"gorm.io/gorm"
)

type Lang struct {
	ID uint `gorm:"type:INT(11) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`

	Code    int    `gorm:"not null;type:INT(11);comment:错误码;uniqueIndex:code"`
	Lang    string `gorm:"column:lang;size:16;comment:语言;uniqueIndex:code"`
	Content string `gorm:"column:zh_HK;size:512;comment:内容"`

	CreatedAt time.Time      `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time      `gorm:"not null;default:current_timestamp on update current_timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_deleted_at"`
}
