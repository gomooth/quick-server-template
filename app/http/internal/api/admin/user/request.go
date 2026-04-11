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

func (in *createRequest) GetGenres() []int8 {
	res := make([]int8, 0)
	for _, genre := range in.Genres {
		if genre == 0 {
			continue
		}
		res = append(res, int8(genre))
	}
	return res
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
	createRequest

	State int8 `json:"state"`
}

func (in *modifyRequest) Validate() error {
	if err := in.createRequest.Validate(); nil != err {
		return err
	}

	return nil
}
