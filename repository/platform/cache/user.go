package cache

import (
	"context"
	"server-api/global"
	"server-api/repository/platform"
	"server-api/repository/platform/dao"
	"server-api/repository/types/platformfilter"
	"time"

	"github.com/gomooth/pkg/framework/dbcache"
	"github.com/gomooth/pkg/framework/dbfilter"

	"github.com/save95/xerror"
)

type user struct {
	name string
}

func NewUser() *user {
	return &user{
		name: "user",
	}
}

func (s *user) getCacher() (dbcache.IDBCache[platform.VWUser, platformfilter.User], error) {
	key := s.name
	cacheManger, err := global.StringCacheManager()
	if err != nil {
		return nil, err
	}

	return dbcache.New[platform.VWUser, platformfilter.User](
		key,
		cacheManger,
		dbcache.WithExpiration(time.Hour), // 默认15分钟，改成1小时缓存
	), nil
}

func (s *user) Paginate(ctx context.Context, start, limit int, opt dbfilter.IFilter[platformfilter.User]) ([]*platform.VWUser, uint, error) {
	cacher, err := s.getCacher()
	if err != nil {
		return nil, 0, err
	}
	return cacher.Paginate(ctx, start, limit, opt, func() ([]*platform.VWUser, uint, error) {
		return dao.NewVWUser().Paginate(ctx, start, limit, opt)
	})
}

func (s *user) First(ctx context.Context, id uint) (*platform.VWUser, error) {
	if id == 0 {
		return nil, xerror.New("id error")
	}
	cacher, err := s.getCacher()
	if err != nil {
		return nil, err
	}
	return cacher.First(ctx, id, func() (*platform.VWUser, error) {
		return dao.NewVWUser().First(ctx, id)
	})
}

func (s *user) Clear(ctx context.Context) error {
	cacher, err := s.getCacher()
	if err != nil {
		return err
	}
	return cacher.Clear(ctx)
}
