package dao

import (
	"context"
	"errors"
	"fmt"
	"server-api/global"
	"server-api/repository/platform"
	"server-api/repository/types/platformfilter"
	"server-api/repository/types/platformtypes"

	"github.com/gomooth/pkg/framework/dbfilter"
	"github.com/gomooth/pkg/framework/dbrepo"

	"github.com/save95/xerror"
	"github.com/save95/xerror/xcode"

	"gorm.io/gorm"
)

type IVWUser interface {
	First(ctx context.Context, id uint, preloads ...string) (*platform.VWUser, error)
	FirstByAccount(ctx context.Context, account string, preloads ...string) (*platform.VWUser, error)
	FirstThirdOpenID(ctx context.Context, third platformtypes.UserFromPlatform, appID string, userID uint) (string, string, error)
	All(ctx context.Context, option dbfilter.IFilter[platformfilter.User]) ([]*platform.VWUser, error)
	List(ctx context.Context, start, limit int, option dbfilter.IFilter[platformfilter.User]) ([]*platform.VWUser, error)
	Paginate(ctx context.Context, start, limit int, option dbfilter.IFilter[platformfilter.User]) ([]*platform.VWUser, uint, error)
	ListRoles(ctx context.Context, id uint) ([]*platform.UserRole, error)
}

type vwUser struct {
	db       *gorm.DB
	dao      dbrepo.IDAO[platform.VWUser]
	searcher dbrepo.ISearcher[platform.VWUser, platformfilter.User]
}

func NewVWUser() IVWUser {
	result := &vwUser{}

	db, _ := global.Database().Get("platform")
	dao := dbrepo.NewDAO[platform.VWUser](db)
	searcher := dbrepo.NewSearcher[platform.VWUser, platformfilter.User](db,
		result.buildFilter, result.getSortKeyMapping(),
	)

	return &vwUser{
		db:       db,
		dao:      dao,
		searcher: searcher,
	}
}

func (u *vwUser) buildFilter(filter *platformfilter.User, db *gorm.DB) *gorm.DB {
	if db == nil {
		db = u.db.Model(platform.VWUser{})
	}

	// 账号
	if v := filter.Account; len(v) > 0 {
		vv := fmt.Sprintf("%%%s%%", v)
		db = db.Where("account like ?", vv)
	}

	return db
}

func (u *vwUser) getSortKeyMapping() map[string]string {
	return map[string]string{
		"created_at":   "created_at",
		"created_time": "created_at",
	}
}

func (u *vwUser) First(ctx context.Context, id uint, preloads ...string) (*platform.VWUser, error) {
	return u.dao.FirstWith(ctx, id, preloads...)
}

func (u *vwUser) FirstByAccount(ctx context.Context, account string, preloads ...string) (*platform.VWUser, error) {
	if len(account) == 0 {
		return nil, xerror.WithXCode(xcode.DBRequestParamError)
	}

	db := u.db.WithContext(ctx).Where("account = ?", account)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	var record platform.VWUser
	if err := db.First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.WithXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}
	return &record, nil
}

// FirstThirdOpenID 获取三方平台用户 unionID/openID
// 返回：unionID, openID, error
func (u *vwUser) FirstThirdOpenID(ctx context.Context, third platformtypes.UserFromPlatform, appID string, userID uint) (
	string /*unionID*/, string /*openID*/, error,
) {
	switch third {
	case platformtypes.UserFromPlatformWechat:
		var record platform.UserWechat
		if err := u.db.WithContext(ctx).Model(platform.UserWechat{}).
			Where("app_id = ?", appID).
			Where("user_id = ?", userID).
			First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", "", xerror.WrapWithXCodeStatus(err, xcode.DBRecordNotFound)
			}
			return "", "", xerror.WrapWithXCodeStatus(err, xcode.DBFailed)
		}
		return record.UnionID, record.OpenID, nil
	case platformtypes.UserFromPlatformAlipay:
		var record platform.UserAlipay
		if err := u.db.WithContext(ctx).Model(platform.UserAlipay{}).
			Where("app_id = ?", appID).
			Where("user_id = ?", userID).
			First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", "", xerror.WrapWithXCodeStatus(err, xcode.DBRecordNotFound)
			}
			return "", "", xerror.WrapWithXCodeStatus(err, xcode.DBFailed)
		}
		return record.UnionID, record.OpenID, nil
	default:
		return "", "", xerror.New("not supported")
	}
}

func (u *vwUser) All(ctx context.Context, option dbfilter.IFilter[platformfilter.User]) ([]*platform.VWUser, error) {
	return u.searcher.All(ctx, option)
}

func (u *vwUser) List(ctx context.Context, start, limit int, option dbfilter.IFilter[platformfilter.User]) ([]*platform.VWUser, error) {
	return u.searcher.List(ctx, start, limit, option)
}

func (u *vwUser) Paginate(ctx context.Context, start, limit int, option dbfilter.IFilter[platformfilter.User]) ([]*platform.VWUser, uint, error) {
	return u.searcher.Paginate(ctx, start, limit, option)
}

func (u *vwUser) ListRoles(ctx context.Context, id uint) ([]*platform.UserRole, error) {
	db := u.db.WithContext(ctx).Model(platform.UserRole{}).
		Where("user_id = ?", id)

	var records []*platform.UserRole
	if err := db.Order("id ASC").Find(&records).Error; err != nil {
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return records, nil
}
