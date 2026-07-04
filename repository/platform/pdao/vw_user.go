package pdao

import (
	"context"
	"errors"
	"fmt"
	"server-api/global"
	"server-api/repository/platform/pattr"
	"server-api/repository/platform/pfilter"
	"server-api/repository/platform/pmodel"

	"github.com/gomooth/pkg/framework/dbquery"
	"github.com/gomooth/pkg/framework/dbrepo"

	"github.com/gomooth/xerror"
	"github.com/gomooth/xerror/xcode"

	"gorm.io/gorm"
)

type IVWUser interface {
	First(ctx context.Context, id uint, preloads ...string) (*pmodel.VWUser, error)
	FirstByAccount(ctx context.Context, account string, preloads ...string) (*pmodel.VWUser, error)
	FirstThirdOpenID(ctx context.Context, third pattr.UserFromPlatform, appID string, userID uint) (string, string, error)
	All(ctx context.Context, option dbquery.IQuery[pfilter.User]) ([]*pmodel.VWUser, error)
	List(ctx context.Context, option dbquery.IQuery[pfilter.User]) ([]*pmodel.VWUser, error)
	Paginate(ctx context.Context, option dbquery.IQuery[pfilter.User]) ([]*pmodel.VWUser, uint, error)
	ListRoles(ctx context.Context, id uint) ([]*pmodel.UserRole, error)
}

type vwUser struct {
	db       *gorm.DB
	dao      dbrepo.IDAO[pmodel.VWUser]
	searcher dbrepo.ISearcher[pmodel.VWUser, pfilter.User]
}

func NewVWUser() IVWUser {
	result := &vwUser{}

	db, _ := global.Database().Get("platform")
	dao, _ := dbrepo.NewDAO[pmodel.VWUser](db)
	searcher, _ := dbrepo.NewSearcher[pmodel.VWUser, pfilter.User](db,
		dbrepo.WithFilterTransfer[pmodel.VWUser, pfilter.User](result.buildFilter),
		dbrepo.WithSortMapping[pmodel.VWUser, pfilter.User](result.getSortMapping()),
	)

	return &vwUser{
		db:       db,
		dao:      dao,
		searcher: searcher,
	}
}

func (u *vwUser) buildFilter(filter *pfilter.User, db *gorm.DB) *gorm.DB {
	if db == nil {
		db = u.db.Model(pmodel.VWUser{})
	}

	// 账号
	if v := filter.Account; len(v) > 0 {
		db = db.Where("account = ?", v)
	}

	// 账号
	if v := filter.AccountLike; len(v) > 0 {
		vv := fmt.Sprintf("%%%s%%", v)
		db = db.Where("account like ?", vv)
	}

	return db
}

func (u *vwUser) getSortMapping() *dbquery.SortMapping {
	return dbquery.NewSortMapping(
		dbquery.WithSortKeyMap(map[string]string{
			"created_at":   "created_at",
			"created_time": "created_at",
		}),
		dbquery.WithDefaultSort("created_at"),
	)
}

func (u *vwUser) First(ctx context.Context, id uint, preloads ...string) (*pmodel.VWUser, error) {
	return u.dao.FirstWith(ctx, id, preloads...)
}

func (u *vwUser) FirstByAccount(ctx context.Context, account string, preloads ...string) (*pmodel.VWUser, error) {
	if len(account) == 0 {
		return nil, xerror.NewXCode(xcode.DBRequestParamError)
	}

	db := u.db.WithContext(ctx).Where("account = ?", account)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	var record pmodel.VWUser
	if err := db.First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.NewXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}
	return &record, nil
}

// FirstThirdOpenID 获取三方平台用户 unionID/openID
// 返回：unionID, openID, error
func (u *vwUser) FirstThirdOpenID(ctx context.Context, third pattr.UserFromPlatform, appID string, userID uint) (
	string /*unionID*/, string /*openID*/, error,
) {
	switch third {
	case pattr.UserFromPlatformWechat:
		var record pmodel.UserWechat
		if err := u.db.WithContext(ctx).Model(pmodel.UserWechat{}).
			Where("app_id = ?", appID).
			Where("user_id = ?", userID).
			First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", "", xerror.WrapStatus(err, xcode.DBRecordNotFound)
			}
			return "", "", xerror.WrapStatus(err, xcode.DBFailed)
		}
		return record.UnionID, record.OpenID, nil
	case pattr.UserFromPlatformAlipay:
		var record pmodel.UserAlipay
		if err := u.db.WithContext(ctx).Model(pmodel.UserAlipay{}).
			Where("app_id = ?", appID).
			Where("user_id = ?", userID).
			First(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", "", xerror.WrapStatus(err, xcode.DBRecordNotFound)
			}
			return "", "", xerror.WrapStatus(err, xcode.DBFailed)
		}
		return record.UnionID, record.OpenID, nil
	default:
		return "", "", xerror.New("not supported")
	}
}

func (u *vwUser) All(ctx context.Context, option dbquery.IQuery[pfilter.User]) ([]*pmodel.VWUser, error) {
	return u.searcher.FindAll(ctx, option)
}

func (u *vwUser) List(ctx context.Context, option dbquery.IQuery[pfilter.User]) ([]*pmodel.VWUser, error) {
	return u.searcher.List(ctx, option)
}

func (u *vwUser) Paginate(ctx context.Context, option dbquery.IQuery[pfilter.User]) ([]*pmodel.VWUser, uint, error) {
	return u.searcher.Paginate(ctx, option)
}

func (u *vwUser) ListRoles(ctx context.Context, id uint) ([]*pmodel.UserRole, error) {
	db := u.db.WithContext(ctx).Model(pmodel.UserRole{}).
		Where("user_id = ?", id)

	var records []*pmodel.UserRole
	if err := db.Order("id ASC").Find(&records).Error; err != nil {
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return records, nil
}
