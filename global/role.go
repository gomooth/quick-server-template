package global

import (
	"errors"

	"github.com/gomooth/pkg/http/httpcontext"
)

type Role uint8

const (
	RoleUser  Role = iota // 用户
	RoleSuper             // 超级管理员
)

var rolesTitle = map[Role]string{
	RoleUser:  "user",
	RoleSuper: "super",
}

func (r Role) String() string {
	if title, ok := rolesTitle[r]; ok {
		return title
	}

	return "unknown"
}

func NewRole(str string) (httpcontext.IRole, error) {
	for i := range rolesTitle {
		if rolesTitle[i] == str {
			return i, nil
		}
	}

	return nil, errors.New("unknown role")
}

func NewRoleFromGenre(genre uint8) (httpcontext.IRole, error) {
	for i := range rolesTitle {
		if uint8(i) == genre {
			return i, nil
		}
	}

	return nil, errors.New("unknown role")
}

func Role2Genre(role httpcontext.IRole) uint8 {
	switch role {
	case RoleUser:
		return uint8(RoleUser)
	case RoleSuper:
		return uint8(RoleSuper)
	default:
		return 0
	}
}
