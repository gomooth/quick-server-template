package dao

import (
	"context"
	"errors"
	"server-api/global"
	"server-api/repository/platform"
	"server-api/repository/types/platformfilter"

	"github.com/gomooth/pkg/framework/dbfilter"
	"github.com/gomooth/pkg/framework/dbrepo"
	"github.com/save95/xerror"
	"github.com/save95/xerror/xcode"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IOpenApp interface {
	Create(ctx context.Context, record *platform.OpenAPP) error
	Creates(ctx context.Context, records []*platform.OpenAPP) error
	CreateOrUpdate(ctx context.Context, record *platform.OpenAPP) error
	CreateOrUpdates(ctx context.Context, records []*platform.OpenAPP) error
	Save(ctx context.Context, record *platform.OpenAPP) error
	Delete(ctx context.Context, id uint) error
	Remove(ctx context.Context, id uint) error
	First(ctx context.Context, id uint) (*platform.OpenAPP, error)
	FirstByAppID(ctx context.Context, appID string) (*platform.OpenAPP, error)
	FirstByOrganization(ctx context.Context, organizationID uint) (*platform.OpenAPP, error)
	All(ctx context.Context, option dbfilter.IFilter[platformfilter.OpenAPP]) ([]*platform.OpenAPP, error)
	List(ctx context.Context, start, limit int, option dbfilter.IFilter[platformfilter.OpenAPP]) ([]*platform.OpenAPP, error)
	Paginate(ctx context.Context, start, limit int, option dbfilter.IFilter[platformfilter.OpenAPP]) ([]*platform.OpenAPP, uint, error)
}

type openApp struct {
	db     *gorm.DB
	dao    dbrepo.IDAO[platform.OpenAPP]
	search dbrepo.ISearcher[platform.OpenAPP, platformfilter.OpenAPP]
}

func NewOpenAPP() IOpenApp {
	result := &openApp{}

	db, _ := global.Database().Get("platform")
	dao := dbrepo.NewDAO[platform.OpenAPP](db)
	queryBuilder := dbrepo.NewSearcher[platform.OpenAPP, platformfilter.OpenAPP](db,
		result.buildFilter, result.getSortKeyMapping())

	return &openApp{
		dao:    dao,
		search: queryBuilder,
	}
}

func (u *openApp) buildFilter(filter *platformfilter.OpenAPP, db *gorm.DB) *gorm.DB {
	if db == nil {
		db = u.db.Model(platform.OpenAPP{})
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

func (u *openApp) getSortKeyMapping() map[string]string {
	return map[string]string{
		"created_at":   "created_at",
		"created_time": "created_at",
	}
}

func (u *openApp) Create(ctx context.Context, record *platform.OpenAPP) error {
	return u.dao.Create(ctx, record)
}

func (u *openApp) Creates(ctx context.Context, records []*platform.OpenAPP) error {
	return u.dao.Creates(ctx, records)
}

func (u *openApp) CreateOrUpdate(ctx context.Context, record *platform.OpenAPP) error {
	return u.CreateOrUpdates(ctx, []*platform.OpenAPP{record})
}

func (u *openApp) CreateOrUpdates(ctx context.Context, records []*platform.OpenAPP) error {
	// 验证参数
	if len(records) == 0 {
		return xerror.WithXCode(xcode.DBRequestParamError)
	}
	for _, record := range records {
		if record.ID > 0 || len(record.AppID) == 0 || len(record.AppSecret) == 0 {
			return xerror.WithXCode(xcode.DBRequestParamError)
		}
	}

	// 使用冲突更新
	if err := u.db.WithContext(ctx).Model(&platform.OpenAPP{}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "app_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"app_secret", "state", "push_url", "push_qps_limit",
			}),
		}).CreateInBatches(records, 100).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return xerror.WithXCode(xcode.DBRecordNotFound)
		}
		return xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return nil
}

func (u *openApp) Save(ctx context.Context, record *platform.OpenAPP) error {
	return u.dao.Save(ctx, record)
}

func (u *openApp) Delete(ctx context.Context, id uint) error {
	return u.dao.Delete(ctx, id)
}

func (u *openApp) Remove(ctx context.Context, id uint) error {
	return u.dao.Remove(ctx, id)
}

func (u *openApp) First(ctx context.Context, id uint) (*platform.OpenAPP, error) {
	return u.dao.First(ctx, id)
}

func (u *openApp) FirstByAppID(ctx context.Context, appID string) (*platform.OpenAPP, error) {
	if len(appID) == 0 {
		return nil, xerror.WithXCode(xcode.DBRequestParamError)
	}

	var record platform.OpenAPP
	if err := u.db.WithContext(ctx).Where("app_id = ?", appID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.WithXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}
	return &record, nil
}

func (u *openApp) FirstByOrganization(ctx context.Context, organizationID uint) (*platform.OpenAPP, error) {
	if organizationID == 0 {
		return nil, xerror.WithXCode(xcode.DBRequestParamError)
	}

	var record platform.OpenAPP
	if err := u.db.WithContext(ctx).Where("organization_id = ?", organizationID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.WithXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}
	return &record, nil
}

func (u *openApp) All(ctx context.Context, option dbfilter.IFilter[platformfilter.OpenAPP]) ([]*platform.OpenAPP, error) {
	return u.search.All(ctx, option)
}

func (u *openApp) List(ctx context.Context, start, limit int, option dbfilter.IFilter[platformfilter.OpenAPP]) ([]*platform.OpenAPP, error) {
	return u.search.List(ctx, start, limit, option)
}

func (u *openApp) Paginate(ctx context.Context, start, limit int, option dbfilter.IFilter[platformfilter.OpenAPP]) ([]*platform.OpenAPP, uint, error) {
	return u.search.Paginate(ctx, start, limit, option)
}
