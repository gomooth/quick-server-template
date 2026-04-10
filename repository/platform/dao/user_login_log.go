package dao

import (
	"context"
	"errors"
	"server-api/global"
	"server-api/repository/platform"
	"server-api/repository/types/platformfilter"

	"github.com/gomooth/pkg/framework/dbfilter"
	"github.com/gomooth/pkg/framework/dbrepo"
	"gorm.io/gorm"

	"github.com/save95/xerror"
	"github.com/save95/xerror/xcode"
)

type IUserLoginLog interface {
	Create(ctx context.Context, record *platform.UserLoginLog) error
	Save(ctx context.Context, record *platform.UserLoginLog) error
	First(ctx context.Context, id uint) (*platform.UserLoginLog, error)
	FirstByUser(ctx context.Context, userID uint) (*platform.UserLoginLog, error)
	Paginate(ctx context.Context, start, limit int, option dbfilter.IFilter[platformfilter.UserLoginLog]) ([]*platform.UserLoginLog, uint, error)
}

type userLoginLog struct {
	db       *gorm.DB
	dao      dbrepo.IDAO[platform.UserLoginLog]
	searcher dbrepo.ISearcher[platform.UserLoginLog, platformfilter.UserLoginLog]
}

func NewUserLoginLog() IUserLoginLog {
	result := &userLoginLog{}

	db, _ := global.Database().Get("platform")
	dao := dbrepo.NewDAO[platform.UserLoginLog](db)
	searcher := dbrepo.NewSearcher[platform.UserLoginLog, platformfilter.UserLoginLog](db,
		result.buildFilter, result.getSortKeyMapping(),
	)

	return &userLoginLog{
		dao:      dao,
		searcher: searcher,
	}
}

func (u *userLoginLog) buildFilter(filter *platformfilter.UserLoginLog, db *gorm.DB) *gorm.DB {
	if db == nil {
		db = u.db.Model(platform.UserLoginLog{})
	}

	if v := filter.UserID; v > 0 {
		db = db.Where("user_id in (?)", v)
	}

	return db
}

func (u *userLoginLog) getSortKeyMapping() map[string]string {
	return map[string]string{
		"created_at":   "created_at",
		"created_time": "created_at",
	}
}

func (u *userLoginLog) Create(ctx context.Context, record *platform.UserLoginLog) error {
	return u.dao.Create(ctx, record)
}

func (u *userLoginLog) Save(ctx context.Context, record *platform.UserLoginLog) error {
	return u.dao.Save(ctx, record)
}

func (u *userLoginLog) First(ctx context.Context, id uint) (*platform.UserLoginLog, error) {
	return u.dao.First(ctx, id)
}

func (u *userLoginLog) FirstByUser(ctx context.Context, userID uint) (*platform.UserLoginLog, error) {
	if userID == 0 {
		return nil, xerror.WithXCode(xcode.DBRequestParamError)
	}

	var record platform.UserLoginLog
	if err := u.db.WithContext(ctx).Where("user_id = ?", userID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.WithXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}
	return &record, nil
}

func (u *userLoginLog) Paginate(ctx context.Context, start, limit int, option dbfilter.IFilter[platformfilter.UserLoginLog]) ([]*platform.UserLoginLog, uint, error) {
	return u.searcher.Paginate(ctx, start, limit, option)
}
