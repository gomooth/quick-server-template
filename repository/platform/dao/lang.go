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
)

type ILang interface {
	FirstByCode(ctx context.Context, code int) (*platform.Lang, error)
	All(ctx context.Context, option dbfilter.IFilter[platformfilter.Lang]) ([]*platform.Lang, error)
	List(ctx context.Context, start, limit int, option dbfilter.IFilter[platformfilter.Lang]) ([]*platform.Lang, error)
	Paginate(ctx context.Context, start, limit int, option dbfilter.IFilter[platformfilter.Lang]) ([]*platform.Lang, uint, error)
}

type lang struct {
	db       *gorm.DB
	dao      dbrepo.IDAO[platform.Lang]
	searcher dbrepo.ISearcher[platform.Lang, platformfilter.Lang]
}

func NewLang() ILang {
	result := &lang{}

	db, _ := global.Database().Get("platform")
	dao := dbrepo.NewDAO[platform.Lang](db)
	searcher := dbrepo.NewSearcher[platform.Lang, platformfilter.Lang](db,
		result.buildFilter, result.getSortKeyMapping(),
	)

	return &lang{
		dao:      dao,
		searcher: searcher,
	}
}

func (l *lang) buildFilter(filter *platformfilter.Lang, db *gorm.DB) *gorm.DB {
	if db == nil {
		db = l.db.Model(platform.Lang{})
	}

	if v := filter.Codes; len(v) > 0 {
		db = db.Where("code in (?)", v)
	}

	if v := filter.Code; v != nil {
		db = db.Where("code = ?", v)
	}

	return db
}

func (l *lang) getSortKeyMapping() map[string]string {
	return map[string]string{
		"created_at":   "created_at",
		"created_time": "created_at",
	}
}

func (l *lang) FirstByCode(ctx context.Context, code int) (*platform.Lang, error) {
	var record platform.Lang
	if err := l.db.WithContext(ctx).Where("code = ?", code).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.WithXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}
	return &record, nil
}

func (l *lang) All(ctx context.Context, option dbfilter.IFilter[platformfilter.Lang]) ([]*platform.Lang, error) {
	return l.searcher.All(ctx, option)
}

func (l *lang) List(ctx context.Context, start, limit int, option dbfilter.IFilter[platformfilter.Lang]) ([]*platform.Lang, error) {
	return l.searcher.List(ctx, start, limit, option)
}

func (l *lang) Paginate(ctx context.Context, start, limit int, option dbfilter.IFilter[platformfilter.Lang]) ([]*platform.Lang, uint, error) {
	return l.searcher.Paginate(ctx, start, limit, option)
}
