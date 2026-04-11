package pmodel

import (
	"time"

	"gorm.io/gorm"
)

type UserAlipay struct {
	gorm.Model

	UserID uint
	BindAt *time.Time

	AppID     string
	OpenID    string
	UnionID   string
	Nickname  string
	AvatarURL string

	Mobile string // 手机号
}
