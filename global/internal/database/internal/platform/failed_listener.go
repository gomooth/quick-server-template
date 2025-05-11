package platform

import (
	"time"

	"gorm.io/gorm"
)

type FailedListener struct {
	ID uint `gorm:"type:INT(11) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`

	ConsumeGroup  string `gorm:"not null;size:64"`
	Topic         string `gorm:"not null;size:32"`
	Msg           string `gorm:"not null;type:MEDIUMTEXT;comment:错误消息"`
	FailedPayload string `gorm:"not null;size:1024;comment:失败原因载荷"`
	FailedReason  string `gorm:"not null;type:MEDIUMTEXT;comment:失败原因"`

	CreatedAt time.Time      `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time      `gorm:"not null;default:current_timestamp on update current_timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_deleted_at"`
}
