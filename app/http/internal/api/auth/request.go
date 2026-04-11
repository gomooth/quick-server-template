package auth

import (
	"server-api/global"

	"github.com/gomooth/xerror"
)

type createTokenRequest struct {
	Genre    global.Role `json:"genre"`
	Account  string      `json:"account"`
	Password string      `json:"password"`
}

func (in createTokenRequest) Validate() error {
	if in.Genre != global.RoleSuper && in.Genre != global.RoleUser {
		return xerror.New("该帐号角色未开通登录权限")
	}

	if len(in.Account) == 0 || len(in.Password) == 0 {
		return xerror.New("请填写正确的登录信息")
	}

	if len(in.Password) < 6 {
		return xerror.New("密码 错误")
	}

	return nil
}

type changePwdRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (in changePwdRequest) Validate() error {
	if len(in.OldPassword) == 0 {
		return xerror.New("原密码 不能为空")
	}
	if len(in.NewPassword) < 6 {
		return xerror.New("新密码 不能少于6个字符")
	}
	if in.NewPassword == in.OldPassword {
		return xerror.New("原密码 和 新密码 不能一致")
	}

	return nil
}
