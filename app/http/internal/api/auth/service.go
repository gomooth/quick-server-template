package auth

import (
	"context"
	"server-api/global"
	"server-api/global/ecode"
	"server-api/app/http/internal/helper"
	"server-api/repository/platform/pdao"
	"server-api/repository/platform/pmodel"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gomooth/pkg/http/httpcontext"
	"github.com/gomooth/pkg/http/jwt"
	"github.com/gomooth/pkg/http/jwt/jwtstore"
	"github.com/gomooth/utils/userutil"

	"github.com/gomooth/xerror"
	"github.com/gomooth/xerror/xcode"
)

type service struct {
}

func (s *service) Login(ctx context.Context, in *createTokenRequest) (*tokenEntity, error) {
	if err := in.Validate(); nil != err {
		return nil, xerror.NewXCode(xcode.RequestParamError, err.Error())
	}

	user, err := pdao.NewVWUser().FirstByAccount(ctx, in.Account, "UserRoles")
	if err != nil {
		return nil, xerror.WrapWithXCode(err, ecode.ErrorAuthParams)
	}

	// 账号无效
	if user.State != 1 {
		return nil, xerror.NewXCode(ecode.ErrorAuthParams)
	}

	// 检查密码
	if !userutil.Check(in.Password, user.Password) {
		return nil, xerror.NewXCode(ecode.ErrorAuthParams)
	}

	return s.makeToken(ctx, user)
}

func (s *service) makeToken(ctx context.Context, user *pmodel.VWUser) (*tokenEntity, error) {

	roles, roleTitles, err := user.Roles()
	if nil != err {
		return nil, err
	}

	// token 有效时长
	duration := 1 * 24 * time.Hour
	// 多地登陆
	store := jwtstore.NewMultiRedisStore(global.SessionStoreClient)
	//// 单一登陆
	//store := jwtstore.NewSingleRedisStore(global.SessionStoreClient)
	// 生成JWT TOKEN
	tk, err := jwt.NewTokenBuilder(
		[]byte(global.Config.App.Secret),
		httpcontext.User{
			ID:      user.ID,
			Account: user.Account,
			Name:    user.Nickname,
			Roles:   roles,
		},
	).WithIssuer(global.Config.App.ID).
		WithExpiration(duration).
		WithStatefulStore(store).
		Build()
	if err != nil {
		return nil, xerror.WrapWithXCode(err, ecode.ErrorAuthFailed)
	}

	tokenStr, err := tk.ToString(ctx)
	if err != nil {
		return nil, xerror.WrapWithXCode(err, ecode.ErrorAuthFailed)
	}

	// 更新最后登陆时间
	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = ctx.(*gin.Context).ClientIP()
	user.UpdatedAt = now
	if err := pdao.NewUser().Save(ctx, user.ToUser()); nil != err {
		return nil, xerror.WrapWithXCode(err, ecode.ErrorAuthFailed)
	}

	// 写登陆日志
	httpRequest := ctx.(*gin.Context).Request
	header := helper.ParseAPPHeader(ctx)
	_ = pdao.NewUserLoginLog().Create(ctx, &pmodel.UserLoginLog{
		UserID:    user.ID,
		UserAgent: httpRequest.UserAgent(),
		IP:        user.LastLoginIP,
		Referer:   httpRequest.Referer(),

		UTMSource:   header.UTMSource(),
		UTMMedium:   header.UTMMedium(),
		UTMCampaign: header.UTMCampaign(),
		UTMTerm:     header.UTMTerm(),
		UTMContent:  header.UTMContent(),
	})

	return &tokenEntity{
		AccessToken: tokenStr,
		ExpireTime:  int64(duration.Seconds()),
		Profile: &profileEntity{
			ID:          user.ID,
			Name:        user.ShowName(),
			AvatarURL:   user.ShowAvatarURL(),
			CurrentRole: user.CurrentRole().String(),
			Roles:       roleTitles,
		},
	}, nil
}

func (s *service) Logout(ctx context.Context) error {
	owner := helper.ParseUser(ctx)
	if owner.GetID() == 0 {
		return nil
	}

	// 清除 token
	if err := jwtstore.NewMultiRedisStore(global.SessionStoreClient).Clean(ctx, owner.GetID()); nil != err {
		return xerror.WrapWithXCode(err, ecode.ErrorHandleFailed)
	}

	return nil
}

func (s *service) ChangePwd(ctx context.Context, in *changePwdRequest) error {
	if err := in.Validate(); nil != err {
		return xerror.NewXCode(xcode.RequestParamError, err.Error())
	}

	owner := helper.ParseUser(ctx)
	if owner.GetID() == 0 {
		return xerror.NewXCode(xcode.Unauthorized)
	}

	user, err := pdao.NewVWUser().First(ctx, owner.ID)
	if nil != err {
		return xerror.WrapWithXCode(err, ecode.ErrorRequestData)
	}

	if !userutil.Check(in.OldPassword, user.Password) {
		return xerror.New("原密码 错误")
	}

	password, err := userutil.Sum(in.NewPassword)
	if nil != err {
		return xerror.WrapWithXCode(err, ecode.ErrorHandleFailed)
	}
	user.Password = password

	if err := pdao.NewUser().Save(ctx, user.ToUser()); nil != err {
		return xerror.WrapWithXCode(err, ecode.ErrorHandleFailed)
	}

	// 注销登陆
	_ = s.Logout(ctx)
	return nil
}
