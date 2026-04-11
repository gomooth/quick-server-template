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
)

type ILang interface {
	FirstByCode(ctx context.Context, code int) (*pmodel.Lang, error)
	All(ctx context.Context, option dbquery.IQuery[pfilter.Lang]) ([]*pmodel.Lang, error)
	List(ctx context.Context, option dbquery.IQuery[pfilter.Lang]) ([]*pmodel.Lang, error)
	Paginate(ctx context.Context, option dbquery.IQuery[pfilter.Lang]) ([]*pmodel.Lang, uint, error)
}

type lang struct {
	db       *gorm.DB
	dao      dbrepo.IDAO[pmodel.Lang]
	searcher dbrepo.ISearcher[pmodel.Lang, pfilter.Lang]
}

func NewLang() ILang {
	result := &lang{}

	db, _ := global.Database().Get("platform")
	dao, _ := dbrepo.NewDAO[pmodel.Lang](db)
	searcher, _ := dbrepo.NewSearcher[pmodel.Lang, pfilter.Lang](db,
		dbrepo.WithFilterTransfer[pmodel.Lang, pfilter.Lang](result.buildFilter),
		dbrepo.WithSortMapping[pmodel.Lang, pfilter.Lang](result.getSortMapping()),
	)

	return &lang{
		dao:      dao,
		searcher: searcher,
	}
}

func (l *lang) buildFilter(filter *pfilter.Lang, db *gorm.DB) *gorm.DB {
	if db == nil {
		db = l.db.Model(pmodel.Lang{})
	}

	if v := filter.Codes; len(v) > 0 {
		db = db.Where("code in (?)", v)
	}

	if v := filter.Code; v != nil {
		db = db.Where("code = ?", v)
	}

	return db
}

func (l *lang) getSortMapping() *dbquery.SortMapping {
	return dbquery.NewSortMapping(
		dbquery.WithSortKeyMap(map[string]string{
			"created_at":   "created_at",
			"created_time": "created_at",
		}),
		dbquery.WithDefaultSort("created_at"),
	)
}

func (l *lang) FirstByCode(ctx context.Context, code int) (*pmodel.Lang, error) {
	var record pmodel.Lang
	if err := l.db.WithContext(ctx).Where("code = ?", code).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.NewXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}
	return &record, nil
}

func (l *lang) All(ctx context.Context, option dbquery.IQuery[pfilter.Lang]) ([]*pmodel.Lang, error) {
	return l.searcher.FindAll(ctx, option)
}

func (l *lang) List(ctx context.Context, option dbquery.IQuery[pfilter.Lang]) ([]*pmodel.Lang, error) {
	return l.searcher.List(ctx, option)
}

func (l *lang) Paginate(ctx context.Context, option dbquery.IQuery[pfilter.Lang]) ([]*pmodel.Lang, uint, error) {
	return l.searcher.Paginate(ctx, option)
}
