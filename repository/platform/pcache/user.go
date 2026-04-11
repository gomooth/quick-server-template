package pcache

import (
	"context"
	"server-api/global"
	"server-api/repository/platform/pdao"
	"server-api/repository/platform/pfilter"
	"server-api/repository/platform/pmodel"
	"time"

	"github.com/gomooth/pkg/framework/dbcache"
	"github.com/gomooth/pkg/framework/dbquery"

	"github.com/gomooth/xerror"
)

type IUserCache interface {
	Paginate(ctx context.Context, opt dbquery.IQuery[pfilter.User]) ([]*pmodel.VWUser, uint, error)
	First(ctx context.Context, id uint) (*pmodel.VWUser, error)
	Clear(ctx context.Context) error
}

type user struct {
	name string
}

func NewUser() IUserCache {
	return &user{
		name: "user",
	}
}

func (s *user) getCacher() (dbcache.IDBCache[pmodel.VWUser, pfilter.User], error) {
	key := s.name
	cacheManger, err := global.StringCacheManager()
	if err != nil {
		return nil, err
	}

	return dbcache.New[pmodel.VWUser, pfilter.User](
		key,
		cacheManger,
		dbcache.WithExpiration(time.Hour), // 默认15分钟，改成1小时缓存
	), nil
}

func (s *user) Paginate(ctx context.Context, opt dbquery.IQuery[pfilter.User]) ([]*pmodel.VWUser, uint, error) {
	cacher, err := s.getCacher()
	if err != nil {
		return nil, 0, err
	}
	return cacher.Paginate(ctx, opt, func(ctx context.Context) ([]*pmodel.VWUser, uint, error) {
		return pdao.NewVWUser().Paginate(ctx, opt)
	})
}

func (s *user) First(ctx context.Context, id uint) (*pmodel.VWUser, error) {
	if id == 0 {
		return nil, xerror.New("id error")
	}
	cacher, err := s.getCacher()
	if err != nil {
		return nil, err
	}
	return cacher.First(ctx, id, func(ctx context.Context) (*pmodel.VWUser, error) {
		return pdao.NewVWUser().First(ctx, id)
	})
}

func (s *user) Clear(ctx context.Context) error {
	cacher, err := s.getCacher()
	if err != nil {
		return err
	}
	return cacher.Clear(ctx, dbcache.ClearWithAll(true))
}
