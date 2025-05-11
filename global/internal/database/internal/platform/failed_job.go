package platform

import (
	"time"

	"gorm.io/gorm"
)

type FailedJob struct {
	ID uint `gorm:"type:INT(11) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`

	JobName     string     `gorm:"not null;size:32;comment:任务名"`
	JobArgs     string     `gorm:"not null;size:1024;comment:任务参数"`
	Payload     string     `gorm:"not null;size:1024;comment:任务载荷"`
	Errors      string     `gorm:"type:mediumtext;comment:错误内容"`
	Handled     bool       `gorm:"not null;comment:是否已处理"`
	HandledAt   *time.Time `gorm:"comment:处理时间"`
	Compensated bool       `gorm:"not null;comment:数据是否已补偿"`

	CreatedAt time.Time      `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time      `gorm:"not null;default:current_timestamp on update current_timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_deleted_at"`
}
