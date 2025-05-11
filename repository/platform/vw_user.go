package platform

import (
	"fmt"
)

type VWUser struct {
	User
}

func (m VWUser) ToUser() *User {
	return &m.User
}

func (u User) ShowAvatarURL() string {
	if len(u.AvatarURL) > 0 {
		return u.AvatarURL
	}

	return ""
}

func (m VWUser) ShowName() string {
	if len(m.Nickname) > 0 {
		return m.Nickname
	}

	if len(m.Account) > 4 {
		return fmt.Sprintf("尾号%s的用户", m.Account[len(m.Account)-4:])
	}

	return "尊贵的用户"
}
