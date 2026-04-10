package cache

import (
	"context"
	"fmt"
	"server-api/global"
	"server-api/repository/platform"
	"server-api/repository/platform/dao"
	"server-api/repository/types/platformfilter"
	"time"

	"github.com/gomooth/pkg/framework/dbcache"
	"github.com/gomooth/pkg/framework/dbfilter"
	"github.com/save95/xerror"
)

type openAPP struct {
	name string
}

func NewOpenAPP() *openAPP {
	return &openAPP{
		name: "openApp",
	}
}

func (s *openAPP) getCacher() (dbcache.IDBCache[platform.OpenAPP, platformfilter.OpenAPP], error) {
	key := s.name
	cacheManger, err := global.StringCacheManager()
	if err != nil {
		return nil, err
	}

	return dbcache.New[platform.OpenAPP, platformfilter.OpenAPP](
		key,
		cacheManger,
		dbcache.WithExpiration(time.Hour), // 默认15分钟，改成1小时缓存
	), nil
}

func (s *openAPP) Paginate(ctx context.Context, start, limit int, opt dbfilter.IFilter[platformfilter.OpenAPP]) ([]*platform.OpenAPP, uint, error) {
	cacher, err := s.getCacher()
	if err != nil {
		return nil, 0, err
	}
	return cacher.Paginate(ctx, start, limit, opt, func() ([]*platform.OpenAPP, uint, error) {
		return dao.NewOpenAPP().Paginate(ctx, start, limit, opt)
	})
}

func (s *openAPP) First(ctx context.Context, id uint) (*platform.OpenAPP, error) {
	if id == 0 {
		return nil, xerror.New("id error")
	}
	cacher, err := s.getCacher()
	if err != nil {
		return nil, err
	}
	return cacher.First(ctx, id, func() (*platform.OpenAPP, error) {
		return dao.NewOpenAPP().First(ctx, id)
	})
}

func (s *openAPP) FirstByAppID(ctx context.Context, appID string) (*platform.OpenAPP, error) {
	if len(appID) == 0 {
		return nil, xerror.New("id error")
	}
	cacher, err := s.getCacher()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("byAppID:%s", appID)
	result, err := cacher.Remember(ctx, key, func() (any, error) {
		return dao.NewOpenAPP().FirstByAppID(ctx, appID)
	})
	if nil != err {
		return nil, err
	}
	return result.(*platform.OpenAPP), nil
}

func (s *openAPP) ClearAll(ctx context.Context) error {
	cacher, err := s.getCacher()
	if err != nil {
		return err
	}
	return cacher.Clear(ctx)
}
