package user

import (
	"context"
	"server-api/global/ecode"
	"server-api/repository/platform/pdao"
	"server-api/repository/platform/pfilter"
	"server-api/repository/platform/pmodel"

	"github.com/gomooth/pkg/framework/dbquery"
	"github.com/gomooth/utils/userutil"

	"github.com/gomooth/xerror"
	"github.com/gomooth/xerror/xcode"

	"github.com/zywaited/xcopy"
)

type service struct {
}

func (s *service) Paginate(ctx context.Context, in *paginateRequest) ([]*entity, uint, error) {
	records, total, err := pdao.NewVWUser().Paginate(ctx, dbquery.NewQuery(
		pfilter.User{
			Account: in.Account,
		},
		dbquery.WithSorts[pfilter.User](in.Sort),
		dbquery.WithOffsetPage[pfilter.User](in.Start, in.Limit),
	))
	if nil != err {
		return nil, 0, err
	}

	var res []*entity
	if err := xcopy.Copy(&res, records); nil != err {
		return nil, 0, xerror.WrapWithXCode(err, ecode.ErrorVOConverted)
	}

	return res, total, nil
}

func (s *service) Create(ctx context.Context, in *createRequest) (*entity, error) {
	if err := in.Validate(); nil != err {
		return nil, xerror.WrapStatus(err, xcode.RequestParamError)
	}

	// 判断重复
	_, err := pdao.NewVWUser().FirstByAccount(ctx, in.Account)
	if nil == err || !xerror.IsXCode(err, xcode.DBRecordNotFound) {
		return nil, xerror.WrapWithXCode(err, ecode.ErrorRecordExist)
	}

	pwd, err := userutil.Sum(in.Password)
	if nil != err {
		return nil, xerror.Wrap(err, "生成密码失败")
	}

	record := pmodel.User{
		Account:   in.Account,
		Nickname:  in.Nickname,
		AvatarURL: in.AvatarURL,
		Password:  pwd,
		State:     1,
	}
	stat := pmodel.UserStat{}
	if err := pdao.NewUser().Create(ctx, in.GetGenres(), &record, &stat); nil != err {
		return nil, xerror.WrapWithXCode(err, ecode.ErrorSavedData)
	}

	var res entity
	if err := xcopy.Copy(&res, record); nil != err {
		return nil, xerror.WrapWithXCode(err, ecode.ErrorVOConverted)
	}

	return &res, nil
}

func (s *service) Modify(ctx context.Context, id uint, in *modifyRequest) (*entity, error) {
	if id == 0 {
		return nil, xerror.NewXCode(ecode.ErrorBadRequest)
	}
	if err := in.Validate(); nil != err {
		return nil, xerror.WrapStatus(err, xcode.RequestParamError)
	}

	record, err := pdao.NewVWUser().First(ctx, id, "UserRoles")
	if nil != err {
		return nil, xerror.WrapWithXCode(err, ecode.ErrorRequestData)
	}

	roles := make([]int8, 0)
	for _, role := range record.UserRoles {
		roles = append(roles, int8(role.Genre))
	}

	// 修改名称
	if record.Account != in.Account {
		// 判断重复
		_, err := pdao.NewVWUser().FirstByAccount(ctx, in.Account)
		if nil == err || !xerror.IsXCode(err, xcode.DBRecordNotFound) {
			return nil, xerror.WrapWithXCode(err, ecode.ErrorRecordExist)
		}

		record.Account = in.Account
	}

	record.Nickname = in.Nickname
	record.AvatarURL = in.AvatarURL
	record.State = in.State

	// 如果密码不为空，则修改密码
	if len(in.Password) != 0 {
		pwd, err := userutil.Sum(in.Password)
		if nil != err {
			return nil, xerror.Wrap(err, "生成密码失败")
		}
		record.Password = pwd
	}

	if err := pdao.NewUser().Update(ctx, record.ToUser(), roles); nil != err {
		return nil, xerror.WrapWithXCode(err, ecode.ErrorSavedData)
	}

	var res entity
	if err := xcopy.Copy(&res, record); nil != err {
		return nil, xerror.WrapWithXCode(err, ecode.ErrorVOConverted)
	}

	return &res, nil
}
