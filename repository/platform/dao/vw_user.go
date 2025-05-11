package dao

import (
	"errors"
	"fmt"
	"server-api/global"
	"server-api/repository/platform"
	"server-api/repository/types/platformfilter"
	"server-api/repository/types/platformtypes"

	"github.com/gomooth/pkg/framework/dbquery"

	"github.com/save95/xerror"
	"github.com/save95/xerror/xcode"

	"gorm.io/gorm"
)

type vwUser struct {
	db *gorm.DB
}

func NewVWUser(options ...interface{}) *vwUser {
	impl := vwUser{}
	for _, option := range options {
		if db, ok := option.(*gorm.DB); ok {
			impl.db = db
		}
	}

	if impl.db == nil {
		impl.db, _ = global.Database().Get("platform")
	}

	return &impl
}

func (u *vwUser) First(id uint, preloads ...string) (*platform.VWUser, error) {
	if id == 0 {
		return nil, xerror.WithXCode(xcode.DBRequestParamError)
	}

	db := u.db.Where("id = ?", id)

	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	var record platform.VWUser
	if err := db.First(&record).Error; nil != err {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.WithXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return &record, nil
}

func (u *vwUser) FirstByAccount(account string, preloads ...string) (*platform.VWUser, error) {
	if len(account) == 0 {
		return nil, xerror.WithXCode(xcode.DBRequestParamError)
	}

	db := u.db.Model(platform.VWUser{}).
		Where("account = ?", account)

	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	var record platform.VWUser
	if err := db.First(&record).Error; nil != err {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.WithXCode(xcode.DBRecordNotFound)
		}

		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return &record, nil
}

// FirstThirdOpenID 获取三方平台用户 unionID/openID
// 返回：unionID, openID, error
func (u *vwUser) FirstThirdOpenID(third platformtypes.UserFromPlatform, appID string, userID uint) (
	string /*unionID*/, string /*openID*/, error,
) {
	switch third {
	case platformtypes.UserFromPlatformWechat:
		var record platform.UserWechat
		if err := u.db.Model(platform.UserWechat{}).
			Where("app_id = ?", appID).
			Where("user_id = ?", userID).
			First(&record).Error; nil != err {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", "", xerror.WrapWithXCodeStatus(err, xcode.DBRecordNotFound)
			}
			return "", "", xerror.WrapWithXCodeStatus(err, xcode.DBFailed)
		}
		return record.UnionID, record.OpenID, nil
	case platformtypes.UserFromPlatformAlipay:
		var record platform.UserAlipay
		if err := u.db.Model(platform.UserAlipay{}).
			Where("app_id = ?", appID).
			Where("user_id = ?", userID).
			First(&record).Error; nil != err {
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

func (u *vwUser) All(option dbquery.IFilter[platformfilter.User]) ([]*platform.VWUser, error) {
	db := u.buildFilter(option.Filter())

	db = option.Build(
		db,
		dbquery.BuildWithSortKeyMappings(u.getSortKeyMapping()),
	)

	var records []*platform.VWUser
	if err := db.Find(&records).Error; nil != err {
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return records, nil
}

func (u *vwUser) List(start, limit int, option dbquery.IFilter[platformfilter.User]) ([]*platform.VWUser, error) {
	db := u.buildFilter(option.Filter())

	db = option.Build(
		db,
		dbquery.BuildWithPage(start, limit),
		dbquery.BuildWithSortKeyMappings(u.getSortKeyMapping()),
	)

	var records []*platform.VWUser
	if err := db.Find(&records).Error; nil != err {
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return records, nil
}

func (u *vwUser) Paginate(start, limit int, option dbquery.IFilter[platformfilter.User]) ([]*platform.VWUser, uint, error) {
	db := u.buildFilter(option.Filter())

	var total int64
	_ = db.Count(&total).Error

	db = option.Build(
		db,
		dbquery.BuildWithPage(start, limit),
		dbquery.BuildWithSortKeyMappings(u.getSortKeyMapping()),
	)

	var records []*platform.VWUser
	if err := db.Find(&records).Error; nil != err {
		return nil, 0, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return records, uint(total), nil
}

func (u *vwUser) buildFilter(filter *platformfilter.User) *gorm.DB {
	db := u.db.Model(platform.VWUser{})

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

func (u *vwUser) ListRoles(id uint) ([]*platform.UserRole, error) {
	db := u.db.Model(platform.UserRole{}).
		Where("user_id = ?", id)

	var records []*platform.UserRole
	if err := db.Order("id ASC").Find(&records).Error; nil != err {
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return records, nil
}
