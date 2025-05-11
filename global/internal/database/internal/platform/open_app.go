package platform

import (
	"time"

	"gorm.io/gorm"
)

type OpenAPP struct {
	ID uint `gorm:"type:INT(11) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`

	AppID        string `gorm:"not null;size:32;uniqueIndex:uk_app"`
	AppSecret    string `gorm:"not null;size:128;comment:密钥"`
	State        int8   `gorm:"not null;default:0;comment:状态"`
	PushURL      string `gorm:"not null;size:256;default:'';comment:数据接收地址"`
	PushQPSLimit uint   `gorm:"column:push_qps_limit;not null;default:0;comment:每秒QPS限制：0-不限制"`

	CreatedAt time.Time      `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time      `gorm:"not null;default:current_timestamp on update current_timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_deleted_at"`
}

func (OpenAPP) TableName() string {
	return "open_apps"
}
