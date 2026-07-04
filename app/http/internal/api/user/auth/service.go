package auth

import (
	"context"
	"time"

	"server-api/global"
	"server-api/global/ecode"
	"server-api/app/http/internal/helper"
	"server-api/repository/platform/pdao"
	"server-api/repository/platform/pmodel"

	"github.com/gomooth/pkg/http/httpcontext"
	"github.com/gomooth/pkg/http/jwt"
	"github.com/gomooth/pkg/http/jwt/jwtstore"
	"github.com/gomooth/utils/userutil"

	"github.com/gomooth/xerror"
	"github.com/gomooth/xerror/xcode"
)

type service struct{}

type loginOptions struct {
	clientIP  string
	userAgent string
	referer   string
}

type loginOption func(*loginOptions)

func withClientIP(ip string) loginOption {
	return func(o *loginOptions) { o.clientIP = ip }
}

func withUserAgent(ua string) loginOption {
	return func(o *loginOptions) { o.userAgent = ua }
}

func withReferer(ref string) loginOption {
	return func(o *loginOptions) { o.referer = ref }
}

func (s *service) Login(ctx context.Context, in *createTokenRequest, opts ...loginOption) (*tokenEntity, error) {
	opt := &loginOptions{}
	for _, o := range opts {
		o(opt)
	}

	user, err := pdao.NewVWUser().FirstByAccount(ctx, in.Account, "UserRoles")
	if err != nil {
		return nil, xerror.WrapWithXCode(err, ecode.ErrorAuthParams)
	}

	if user.State != 1 {
		return nil, xerror.NewXCode(ecode.ErrorAuthParams)
	}

	// 必须是 User 角色
	if user.CurrentRole() != global.RoleUser {
		return nil, xerror.NewXCode(ecode.ErrorAuthParams)
	}

	if !userutil.Check(in.Password, user.Password) {
		return nil, xerror.NewXCode(ecode.ErrorAuthParams)
	}

	return s.makeToken(ctx, user, opt)
}

func (s *service) makeToken(ctx context.Context, user *pmodel.VWUser, opt *loginOptions) (*tokenEntity, error) {
	roles, roleTitles, err := user.Roles()
	if nil != err {
		return nil, err
	}

	duration := 1 * 24 * time.Hour
	// 单一登陆
	store := jwtstore.NewSingleRedisStore(global.SessionStoreClient)
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

	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = opt.clientIP
	user.UpdatedAt = now
	if err := pdao.NewUser().Save(ctx, user.ToUser()); nil != err {
		return nil, xerror.WrapWithXCode(err, ecode.ErrorAuthFailed)
	}

	header := helper.ParseAPPHeader(ctx)
	_ = pdao.NewUserLoginLog().Create(ctx, &pmodel.UserLoginLog{
		UserID:    user.ID,
		UserAgent: opt.userAgent,
		IP:        opt.clientIP,
		Referer:   opt.referer,

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

	if err := jwtstore.NewSingleRedisStore(global.SessionStoreClient).Clean(ctx, owner.GetID()); nil != err {
		return xerror.WrapWithXCode(err, ecode.ErrorHandleFailed)
	}

	return nil
}

func (s *service) ChangePwd(ctx context.Context, in *changePwdRequest) error {
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
		return xerror.WrapWithXCode(err, ecode.ErrorPasswordFailed)
	}
	user.Password = password

	if err := pdao.NewUser().Save(ctx, user.ToUser()); nil != err {
		return xerror.WrapWithXCode(err, ecode.ErrorHandleFailed)
	}

	_ = s.Logout(ctx)
	return nil
}
