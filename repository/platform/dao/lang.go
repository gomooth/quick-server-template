package dao

import (
	"errors"
	"server-api/global"
	"server-api/repository/platform"
	"server-api/repository/types/platformfilter"

	"github.com/gomooth/pkg/framework/dbquery"

	"github.com/save95/xerror"
	"github.com/save95/xerror/xcode"

	"gorm.io/gorm"
)

type lang struct {
	db *gorm.DB
}

func NewLang(options ...interface{}) *lang {
	impl := lang{}
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

func (u *lang) First(id uint) (*platform.Lang, error) {
	if id == 0 {
		return nil, xerror.WithXCode(xcode.DBRequestParamError)
	}

	var record platform.Lang
	if err := u.db.Where("id = ?", id).First(&record).Error; nil != err {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.WithXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return &record, nil
}

func (u *lang) FirstByCode(code int) (*platform.Lang, error) {
	db := u.db.Model(platform.Lang{}).
		Where("code = ?", code)

	var record platform.Lang
	if err := db.First(&record).Error; nil != err {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.WithXCode(xcode.DBRecordNotFound)
		}

		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return &record, nil
}

func (u *lang) All(option dbquery.IFilter[platformfilter.Lang]) ([]*platform.Lang, error) {
	db := u.buildFilter(option.Filter())

	db = option.Build(db,
		dbquery.BuildWithSortKeyMappings(u.getSortKeyMapping()),
	)

	var records []*platform.Lang
	if err := db.Find(&records).Error; nil != err {
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return records, nil
}

func (u *lang) List(start, limit int, option dbquery.IFilter[platformfilter.Lang]) ([]*platform.Lang, error) {
	db := u.buildFilter(option.Filter())

	db = option.Build(db,
		dbquery.BuildWithPage(start, limit),
		dbquery.BuildWithSortKeyMappings(u.getSortKeyMapping()),
	)

	var records []*platform.Lang
	if err := db.Find(&records).Error; nil != err {
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return records, nil
}

func (u *lang) Paginate(start, limit int, option dbquery.IFilter[platformfilter.Lang]) ([]*platform.Lang, uint, error) {
	db := u.buildFilter(option.Filter())

	var total int64
	_ = db.Count(&total).Error

	db = option.Build(db,
		dbquery.BuildWithPage(start, limit),
		dbquery.BuildWithSortKeyMappings(u.getSortKeyMapping()),
	)

	var records []*platform.Lang
	if err := db.Find(&records).Error; nil != err {
		return nil, 0, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return records, uint(total), nil
}

func (u *lang) buildFilter(filter *platformfilter.Lang) *gorm.DB {
	db := u.db.Model(platform.Lang{})

	if v := filter.Codes; len(v) > 0 {
		db = db.Where("code in (?)", v)
	}

	if v := filter.Code; v != nil {
		db = db.Where("code = ?", v)
	}

	return db
}

func (u *lang) getSortKeyMapping() map[string]string {
	return map[string]string{
		"created_at":   "created_at",
		"created_time": "created_at",
	}
}
