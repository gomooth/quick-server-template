package user

import (
	"github.com/gomooth/pkg/http/httpmodel"
	"github.com/gomooth/xerror"
)

type paginateRequest struct {
	httpmodel.SearchRequest

	Account string `form:"account"`
}

type createRequest struct {
	Account   string  `json:"account"`
	Nickname  string  `json:"nickname"`
	AvatarURL string  `json:"avatarUrl"`
	Password  string  `json:"password"`
	Genres    []uint8 `json:"genres"`
}

func (in *createRequest) Validate() error {
	if len(in.Account) == 0 || len(in.Genres) == 0 {
		return xerror.New("帐号、类型 不能为空")
	}

	if len(in.Password) == 0 {
		return xerror.New("密码不能为空")
	}

	return nil
}

type modifyRequest struct {
	Account   string  `json:"account"`
	Nickname  string  `json:"nickname"`
	AvatarURL string  `json:"avatarUrl"`
	Password  string  `json:"password"`
	Genres    []uint8 `json:"genres"`
	State     int8    `json:"state"`
}

func (in *modifyRequest) Validate() error {
	if len(in.Account) == 0 {
		return xerror.New("帐号不能为空")
	}
	if in.State != 0 && in.State != 1 {
		return xerror.New("状态值不合法")
	}

	return nil
}
