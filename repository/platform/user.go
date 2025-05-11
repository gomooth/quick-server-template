package platform

import (
	"server-api/global"
	"server-api/repository/types/platformtypes"
	"time"

	"github.com/gomooth/pkg/http/httpcontext"

	"github.com/save95/xerror"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	CheckedGenre uint8
	Account      string
	Nickname     string
	Gender       platformtypes.Gender
	AvatarURL    string
	Password     string
	State        int8
	LastLoginAt  *time.Time
	LastLoginIP  string
	LastDriverNo string

	UserRoles []*UserRole `gorm:"foreignKey:UserID"`
}

func (u User) HasPassword() bool {
	return len(u.Password) > 0
}

func (u User) StateText() string {
	switch u.State {
	case -1:
		return "禁用"
	case 0:
		return "未知"
	case 1:
		return "有效"
	default:
		return "未定义"
	}
}

func (u User) Genres() []int8 {
	res := []int8{
		int8(global.RoleUser),
	}

	if u.UserRoles == nil || len(u.UserRoles) == 0 {
		return res
	}

	for _, role := range u.UserRoles {
		res = append(res, int8(role.Genre))
	}

	return res
}

func (u User) Roles() ([]httpcontext.IRole, []string, error) {
	if u.UserRoles == nil || len(u.UserRoles) == 0 {
		return nil, nil, xerror.New("user roles is empty")
	}

	roles := []httpcontext.IRole{
		global.RoleUser,
	}
	titles := []string{
		global.RoleUser.String(),
	}

	for _, role := range u.UserRoles {
		r, err := global.NewRoleFromGenre(role.Genre)
		if nil != err {
			return nil, nil, err
		}
		roles = append(roles, r)
		titles = append(titles, r.String())
	}

	return roles, titles, nil
}

func (u User) CurrentRole() httpcontext.IRole {
	r, err := global.NewRoleFromGenre(u.CheckedGenre)
	if nil == err {
		return r
	}

	return global.RoleUser
}
