package pdao

import (
	"context"
	"errors"
	"server-api/global"
	"server-api/repository/platform/pfilter"
	"server-api/repository/platform/pmodel"

	"github.com/gomooth/pkg/framework/dbquery"
	"github.com/gomooth/pkg/framework/dbrepo"
	"gorm.io/gorm"

	"github.com/gomooth/xerror"
	"github.com/gomooth/xerror/xcode"
)

type IUserLoginLog interface {
	Create(ctx context.Context, record *pmodel.UserLoginLog) error
	Save(ctx context.Context, record *pmodel.UserLoginLog) error
	First(ctx context.Context, id uint) (*pmodel.UserLoginLog, error)
	FirstByUser(ctx context.Context, userID uint) (*pmodel.UserLoginLog, error)
	Paginate(ctx context.Context, option dbquery.IQuery[pfilter.UserLoginLog]) ([]*pmodel.UserLoginLog, uint, error)
}

type userLoginLog struct {
	db       *gorm.DB
	dao      dbrepo.IDAO[pmodel.UserLoginLog]
	searcher dbrepo.ISearcher[pmodel.UserLoginLog, pfilter.UserLoginLog]
}

func NewUserLoginLog() IUserLoginLog {
	result := &userLoginLog{}

	db, _ := global.Database().Get("platform")
	dao, _ := dbrepo.NewDAO[pmodel.UserLoginLog](db)
	searcher, _ := dbrepo.NewSearcher[pmodel.UserLoginLog, pfilter.UserLoginLog](db,
		dbrepo.WithFilterTransfer[pmodel.UserLoginLog, pfilter.UserLoginLog](result.buildFilter),
		dbrepo.WithSortMapping[pmodel.UserLoginLog, pfilter.UserLoginLog](result.getSortMapping()),
	)

	return &userLoginLog{
		dao:      dao,
		searcher: searcher,
	}
}

func (u *userLoginLog) buildFilter(filter *pfilter.UserLoginLog, db *gorm.DB) *gorm.DB {
	if db == nil {
		db = u.db.Model(pmodel.UserLoginLog{})
	}

	if v := filter.UserID; v > 0 {
		db = db.Where("user_id in (?)", v)
	}

	return db
}

func (u *userLoginLog) getSortMapping() *dbquery.SortMapping {
	return dbquery.NewSortMapping(
		dbquery.WithSortKeyMap(map[string]string{
			"created_at":   "created_at",
			"created_time": "created_at",
		}),
		dbquery.WithDefaultSort("created_at"),
	)
}

func (u *userLoginLog) Create(ctx context.Context, record *pmodel.UserLoginLog) error {
	return u.dao.Create(ctx, record)
}

func (u *userLoginLog) Save(ctx context.Context, record *pmodel.UserLoginLog) error {
	return u.dao.Save(ctx, record)
}

func (u *userLoginLog) First(ctx context.Context, id uint) (*pmodel.UserLoginLog, error) {
	return u.dao.First(ctx, id)
}

func (u *userLoginLog) FirstByUser(ctx context.Context, userID uint) (*pmodel.UserLoginLog, error) {
	if userID == 0 {
		return nil, xerror.NewXCode(xcode.DBRequestParamError)
	}

	var record pmodel.UserLoginLog
	if err := u.db.WithContext(ctx).Where("user_id = ?", userID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.NewXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}
	return &record, nil
}

func (u *userLoginLog) Paginate(ctx context.Context, option dbquery.IQuery[pfilter.UserLoginLog]) ([]*pmodel.UserLoginLog, uint, error) {
	return u.searcher.Paginate(ctx, option)
}
