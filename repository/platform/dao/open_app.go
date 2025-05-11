package dao

import (
	"errors"
	"server-api/global"
	"server-api/repository/platform"
	"server-api/repository/types/platformfilter"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/gomooth/pkg/framework/dbquery"

	"github.com/save95/xerror"
	"github.com/save95/xerror/xcode"
)

type openApp struct {
	db *gorm.DB
}

func NewOpenAPP(options ...interface{}) *openApp {
	impl := openApp{}
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

func (u *openApp) Create(record *platform.OpenAPP) error {
	return u.Creates([]*platform.OpenAPP{record})
}

func (u *openApp) Creates(records []*platform.OpenAPP) error {
	for _, item := range records {
		if item.ID > 0 {
			return xerror.WithXCode(xcode.DBRequestParamError)
		}
	}

	if err := u.db.CreateInBatches(records, 100).Error; nil != err {
		return xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return nil
}

func (u *openApp) CreateOrUpdate(record *platform.OpenAPP) error {
	return u.CreateOrUpdates([]*platform.OpenAPP{record})
}

func (u *openApp) CreateOrUpdates(records []*platform.OpenAPP) error {
	for _, record := range records {
		if record.ID > 0 || len(record.AppID) == 0 || len(record.AppSecret) == 0 {
			return xerror.WithXCode(xcode.DBRequestParamError)
		}
	}

	if err := u.db.Model(&platform.OpenAPP{}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "app_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"app_secret", "state", "push_url", "push_qps_limit",
			}),
		}).CreateInBatches(records, 100).Error; nil != err {
		return xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return nil
}

func (u *openApp) Save(record *platform.OpenAPP) error {
	if record.ID == 0 {
		return xerror.WithXCode(xcode.DBRequestParamError)
	}

	if err := u.db.Save(record).Error; nil != err {
		return xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return nil
}

func (u *openApp) Delete(id uint) error {
	if id == 0 {
		return xerror.WithXCode(xcode.DBRequestParamError)
	}

	var record platform.OpenAPP
	if err := u.db.Where("id = ?", id).Delete(&record).Error; nil != err {
		return xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return nil
}

func (u *openApp) Remove(id uint) error {
	if id == 0 {
		return xerror.WithXCode(xcode.DBRequestParamError)
	}

	var record platform.OpenAPP
	if err := u.db.Unscoped().Where("id = ?", id).Delete(&record).Error; nil != err {
		return xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return nil
}

func (u *openApp) First(id uint) (*platform.OpenAPP, error) {
	if id == 0 {
		return nil, xerror.WithXCode(xcode.DBRequestParamError)
	}

	var record platform.OpenAPP
	if err := u.db.Where("id = ?", id).First(&record).Error; nil != err {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.WithXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return &record, nil
}

func (u *openApp) FirstByAppID(appID string) (*platform.OpenAPP, error) {
	if len(appID) == 0 {
		return nil, xerror.WithXCode(xcode.DBRequestParamError)
	}

	var record platform.OpenAPP
	if err := u.db.Where("app_id = ?", appID).First(&record).Error; nil != err {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.WithXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return &record, nil
}

func (u *openApp) FirstByOrganization(organizationID uint) (*platform.OpenAPP, error) {
	if organizationID == 0 {
		return nil, xerror.WithXCode(xcode.DBRequestParamError)
	}

	var record platform.OpenAPP
	if err := u.db.Where("organization_id = ?", organizationID).First(&record).Error; nil != err {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.WithXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return &record, nil
}

func (u *openApp) All(option dbquery.IFilter[platformfilter.OpenAPP]) ([]*platform.OpenAPP, error) {
	db := u.buildFilter(option.Filter())

	db = option.Build(db)

	var records []*platform.OpenAPP
	if err := db.Find(&records).Error; nil != err {
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return records, nil
}

func (u *openApp) List(start, limit int, option dbquery.IFilter[platformfilter.OpenAPP]) ([]*platform.OpenAPP, error) {
	db := u.buildFilter(option.Filter())

	db = option.Build(
		db,
		dbquery.BuildWithPage(start, limit),
	)

	var records []*platform.OpenAPP
	if err := db.Find(&records).Error; nil != err {
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return records, nil
}

func (u *openApp) Paginate(start, limit int, option dbquery.IFilter[platformfilter.OpenAPP]) ([]*platform.OpenAPP, uint, error) {
	db := u.buildFilter(option.Filter())

	var total int64
	_ = db.Count(&total).Error

	db = option.Build(
		db,
		dbquery.BuildWithPage(start, limit),
	)

	var records []*platform.OpenAPP
	if err := db.Find(&records).Error; nil != err {
		return nil, 0, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return records, uint(total), nil
}

func (u *openApp) buildFilter(filter *platformfilter.OpenAPP) *gorm.DB {
	db := u.db.Model(platform.OpenAPP{})

	if filter.ID > 0 {
		db = db.Where("id = ?", filter.ID)
	}

	if len(filter.IDs) > 0 {
		db = db.Where("id in (?)", filter.IDs)
	}

	if len(filter.AppID) > 0 {
		db = db.Where("app_id = ?", filter.AppID)
	}

	if filter.State != nil {
		db = db.Where("state = ?", filter.State)
	}

	return db
}
