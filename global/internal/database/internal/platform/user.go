package platform

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID uint `gorm:"type:INT(11) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`

	Account      string     `gorm:"not null;size:32;uniqueIndex:uk_account;comment:账号"`
	Nickname     string     `gorm:"not null;size:32;default:'';comment:昵称"`
	CheckedGenre uint8      `gorm:"not null;comment:当前选择角色"`
	Gender       uint8      `gorm:"not null;default:0;comment:性别"`
	AvatarURL    string     `gorm:"not null;size:128;default:'';comment:头像地址"`
	Password     string     `gorm:"not null;size:128;comment:密码"`
	State        int8       `gorm:"not null;default:0;comment:状态"`
	LastLoginAt  *time.Time `gorm:"comment:最后登录时间"`
	LastLoginIP  string     `gorm:"size:32;comment:最后登录IP"`
	LastDriverNo string     `gorm:"not null;size:32;default:'';comment:最后登录设备号"`

	CreatedAt time.Time      `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time      `gorm:"not null;default:current_timestamp on update current_timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index:idx_deleted_at"`
}
