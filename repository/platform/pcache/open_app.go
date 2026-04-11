package pcache

import (
	"context"
	"encoding/json"
	"fmt"
	"server-api/global"
	"server-api/repository/platform/pdao"
	"server-api/repository/platform/pfilter"
	"server-api/repository/platform/pmodel"
	"time"

	"github.com/gomooth/pkg/framework/dbcache"
	"github.com/gomooth/pkg/framework/dbquery"
	"github.com/gomooth/xerror"
)

type IOpenAPPCache interface {
	Paginate(ctx context.Context, opt dbquery.IQuery[pfilter.OpenAPP]) ([]*pmodel.OpenAPP, uint, error)
	First(ctx context.Context, id uint) (*pmodel.OpenAPP, error)
	FirstByAppID(ctx context.Context, appID string) (*pmodel.OpenAPP, error)
	ClearAll(ctx context.Context) error
}

type openAPP struct {
	name string
}

func NewOpenAPP() IOpenAPPCache {
	return &openAPP{
		name: "openApp",
	}
}

func (s *openAPP) getCacher() (dbcache.IDBCache[pmodel.OpenAPP, pfilter.OpenAPP], error) {
	key := s.name
	cacheManger, err := global.StringCacheManager()
	if err != nil {
		return nil, err
	}

	return dbcache.New[pmodel.OpenAPP, pfilter.OpenAPP](
		key,
		cacheManger,
		dbcache.WithExpiration(time.Hour), // 默认15分钟，改成1小时缓存
	), nil
}

func (s *openAPP) Paginate(ctx context.Context, opt dbquery.IQuery[pfilter.OpenAPP]) ([]*pmodel.OpenAPP, uint, error) {
	cacher, err := s.getCacher()
	if err != nil {
		return nil, 0, err
	}
	return cacher.Paginate(ctx, opt, func(ctx context.Context) ([]*pmodel.OpenAPP, uint, error) {
		return pdao.NewOpenAPP().Paginate(ctx, opt)
	})
}

func (s *openAPP) First(ctx context.Context, id uint) (*pmodel.OpenAPP, error) {
	if id == 0 {
		return nil, xerror.New("id error")
	}
	cacher, err := s.getCacher()
	if err != nil {
		return nil, err
	}
	return cacher.First(ctx, id, func(ctx context.Context) (*pmodel.OpenAPP, error) {
		return pdao.NewOpenAPP().First(ctx, id)
	})
}

func (s *openAPP) FirstByAppID(ctx context.Context, appID string) (*pmodel.OpenAPP, error) {
	if len(appID) == 0 {
		return nil, xerror.New("id error")
	}
	cacher, err := s.getCacher()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("byAppID:%s", appID)
	result, err := cacher.Remember(ctx, key, func(ctx context.Context) ([]byte, error) {
		record, err := pdao.NewOpenAPP().FirstByAppID(ctx, appID)
		if err != nil {
			return nil, err
		}
		return json.Marshal(record)
	})
	if err != nil {
		return nil, err
	}
	var record pmodel.OpenAPP
	if err := json.Unmarshal(result, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (s *openAPP) ClearAll(ctx context.Context) error {
	cacher, err := s.getCacher()
	if err != nil {
		return err
	}
	return cacher.Clear(ctx, dbcache.ClearWithAll(true))
}
