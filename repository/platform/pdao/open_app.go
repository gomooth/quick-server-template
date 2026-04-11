package pdao

import (
	"context"
	"errors"
	"server-api/global"
	"server-api/repository/platform/pfilter"
	"server-api/repository/platform/pmodel"

	"github.com/gomooth/pkg/framework/dbquery"
	"github.com/gomooth/pkg/framework/dbrepo"
	"github.com/gomooth/xerror"
	"github.com/gomooth/xerror/xcode"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IOpenApp interface {
	Create(ctx context.Context, record *pmodel.OpenAPP) error
	Creates(ctx context.Context, records []*pmodel.OpenAPP) error
	CreateOrUpdate(ctx context.Context, record *pmodel.OpenAPP) error
	CreateOrUpdates(ctx context.Context, records []*pmodel.OpenAPP) error
	Save(ctx context.Context, record *pmodel.OpenAPP) error
	Delete(ctx context.Context, id uint) (int64, error)
	Remove(ctx context.Context, id uint) (int64, error)
	First(ctx context.Context, id uint) (*pmodel.OpenAPP, error)
	FirstByAppID(ctx context.Context, appID string) (*pmodel.OpenAPP, error)
	FirstByOrganization(ctx context.Context, organizationID uint) (*pmodel.OpenAPP, error)
	All(ctx context.Context, option dbquery.IQuery[pfilter.OpenAPP]) ([]*pmodel.OpenAPP, error)
	List(ctx context.Context, option dbquery.IQuery[pfilter.OpenAPP]) ([]*pmodel.OpenAPP, error)
	Paginate(ctx context.Context, option dbquery.IQuery[pfilter.OpenAPP]) ([]*pmodel.OpenAPP, uint, error)
}

type openApp struct {
	db     *gorm.DB
	dao    dbrepo.IDAO[pmodel.OpenAPP]
	search dbrepo.ISearcher[pmodel.OpenAPP, pfilter.OpenAPP]
}

func NewOpenAPP() IOpenApp {
	result := &openApp{}

	db, _ := global.Database().Get("platform")
	dao, _ := dbrepo.NewDAO[pmodel.OpenAPP](db)
	searcher, _ := dbrepo.NewSearcher[pmodel.OpenAPP, pfilter.OpenAPP](db,
		dbrepo.WithFilterTransfer[pmodel.OpenAPP, pfilter.OpenAPP](result.buildFilter),
		dbrepo.WithSortMapping[pmodel.OpenAPP, pfilter.OpenAPP](result.getSortMapping()),
	)

	return &openApp{
		dao:    dao,
		search: searcher,
	}
}

func (u *openApp) buildFilter(filter *pfilter.OpenAPP, db *gorm.DB) *gorm.DB {
	if db == nil {
		db = u.db.Model(pmodel.OpenAPP{})
	}

	if v := filter.ID; v > 0 {
		db = db.Where("id = ?", v)
	}
	if v := filter.IDs; len(v) > 0 {
		db = db.Where("id in (?)", v)
	}

	if v := filter.AppID; len(v) > 0 {
		db = db.Where("app_id = ?", v)
	}

	if v := filter.State; v != nil {
		db = db.Where("state= ?", v)
	}

	return db
}

func (u *openApp) getSortMapping() *dbquery.SortMapping {
	return dbquery.NewSortMapping(
		dbquery.WithSortKeyMap(map[string]string{
			"created_at":   "created_at",
			"created_time": "created_at",
		}),
		dbquery.WithDefaultSort("created_at"),
	)
}

func (u *openApp) Create(ctx context.Context, record *pmodel.OpenAPP) error {
	return u.dao.Create(ctx, record)
}

func (u *openApp) Creates(ctx context.Context, records []*pmodel.OpenAPP) error {
	return u.dao.Creates(ctx, records)
}

func (u *openApp) CreateOrUpdate(ctx context.Context, record *pmodel.OpenAPP) error {
	return u.CreateOrUpdates(ctx, []*pmodel.OpenAPP{record})
}

func (u *openApp) CreateOrUpdates(ctx context.Context, records []*pmodel.OpenAPP) error {
	// 验证参数
	if len(records) == 0 {
		return xerror.NewXCode(xcode.DBRequestParamError)
	}
	for _, record := range records {
		if record.ID > 0 || len(record.AppID) == 0 || len(record.AppSecret) == 0 {
			return xerror.NewXCode(xcode.DBRequestParamError)
		}
	}

	// 使用冲突更新
	if err := u.db.WithContext(ctx).Model(&pmodel.OpenAPP{}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "app_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"app_secret", "state", "push_url", "push_qps_limit",
			}),
		}).CreateInBatches(records, 100).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return xerror.NewXCode(xcode.DBRecordNotFound)
		}
		return xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return nil
}

func (u *openApp) Save(ctx context.Context, record *pmodel.OpenAPP) error {
	return u.dao.Save(ctx, record)
}

func (u *openApp) Delete(ctx context.Context, id uint) (int64, error) {
	return u.dao.Delete(ctx, id)
}

func (u *openApp) Remove(ctx context.Context, id uint) (int64, error) {
	return u.dao.Remove(ctx, id)
}

func (u *openApp) First(ctx context.Context, id uint) (*pmodel.OpenAPP, error) {
	return u.dao.First(ctx, id)
}

func (u *openApp) FirstByAppID(ctx context.Context, appID string) (*pmodel.OpenAPP, error) {
	if len(appID) == 0 {
		return nil, xerror.NewXCode(xcode.DBRequestParamError)
	}

	var record pmodel.OpenAPP
	if err := u.db.WithContext(ctx).Where("app_id = ?", appID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.NewXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}
	return &record, nil
}

func (u *openApp) FirstByOrganization(ctx context.Context, organizationID uint) (*pmodel.OpenAPP, error) {
	if organizationID == 0 {
		return nil, xerror.NewXCode(xcode.DBRequestParamError)
	}

	var record pmodel.OpenAPP
	if err := u.db.WithContext(ctx).Where("organization_id = ?", organizationID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.NewXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}
	return &record, nil
}

func (u *openApp) All(ctx context.Context, option dbquery.IQuery[pfilter.OpenAPP]) ([]*pmodel.OpenAPP, error) {
	return u.search.FindAll(ctx, option)
}

func (u *openApp) List(ctx context.Context, option dbquery.IQuery[pfilter.OpenAPP]) ([]*pmodel.OpenAPP, error) {
	return u.search.List(ctx, option)
}

func (u *openApp) Paginate(ctx context.Context, option dbquery.IQuery[pfilter.OpenAPP]) ([]*pmodel.OpenAPP, uint, error) {
	return u.search.Paginate(ctx, option)
}
