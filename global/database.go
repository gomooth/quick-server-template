package global

import (
	"context"
	"server-api/global/internal/database"
	inlocker "server-api/global/internal/locker"
	"strings"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	redisstore "github.com/eko/gocache/store/redis/v4"

	"github.com/redis/go-redis/v9"

	"github.com/gomooth/locker"
	"github.com/gomooth/pkg/framework/dbmanager"
	"github.com/gomooth/pkg/framework/dbutil"

	"github.com/gomooth/xerror"
)


var (
	_redisClient  *redis.Client
	_locker       locker.ILocker
	_cacheManager *cache.Cache[string]

	SessionStoreClient *redis.Client
)

func Database() dbmanager.IDatabaseManager {
	return database.Database()
}

// RedisClient redis 实例
func RedisClient() (*redis.Client, error) {
	if _redisClient == nil {
		return nil, xerror.New("redis client is disabled")
	}

	return _redisClient, nil
}

// Locker 分布式锁
func Locker() (locker.ILocker, error) {
	if _locker == nil {
		return nil, xerror.New("locker is disabled")
	}
	return _locker, nil
}

// StringCacheManager 字符串缓存管理器
func StringCacheManager() (*cache.Cache[string], error) {
	if _cacheManager == nil {
		return nil, xerror.New("cache manager is disabled")
	}
	return _cacheManager, nil
}

func InitDataBase() error {
	if err := initDB(); nil != err {
		return xerror.Wrap(err, "database init failed")
	}

	if err := initLocker(); nil != err {
		return xerror.Wrap(err, "locker init failed")
	}

	if err := initRedis(); nil != err {
		return xerror.Wrap(err, "redis init failed")
	}

	if err := initCache(); nil != err {
		return xerror.Wrap(err, "cache init failed")
	}

	return nil
}

func initDB() error {
	if !Config.Data.Persistent.Enabled {
		Log.Debug("database disabled, skip")
		return nil
	}

	for _, db := range Config.Data.Persistent.Connects {
		c, err := dbutil.Connect(&dbutil.Option{
			Name: db.Name,
			Config: &dbutil.ConnectConfig{
				Dsn:         db.Dsn,
				Driver:      db.Driver,
				MaxIdle:     db.MaxIdle,
				MaxOpen:     db.MaxOpen,
				LogMode:     db.LogMode,
				MaxLifeTime: db.MaxLifeTime,
			},
			Logger: Log,
		})
		if err != nil {
			return err
		}

		if err := Database().Register(db.Name, c); nil != err {
			return err
		}
	}

	return nil
}

func initRedis() error {
	if !Config.Data.Redis.Enabled {
		Log.Debug("redis disabled, skip")
		return nil
	}

	_redisClient = redis.NewClient(&redis.Options{
		Addr:     Config.Data.Redis.Addr,
		Password: Config.Data.Redis.Password,
		DB:       Config.Data.Redis.DB,
	})
	if err := _redisClient.Ping(context.Background()).Err(); nil != err {
		return xerror.Wrap(err, "redis client connect failed")
	}

	SessionStoreClient = redis.NewClient(&redis.Options{
		Addr:     Config.Data.Redis.Addr,
		Password: Config.Data.Redis.Password,
		DB:       10,
	})
	if err := SessionStoreClient.Ping(context.Background()).Err(); nil != err {
		return xerror.Wrap(err, "redis client connect failed")
	}

	Log.Debug("redis enabled, init success")

	return nil
}

func initCache() error {
	if !Config.Data.Cache.Enabled {
		Log.Debug("cache disabled, skip")
		return nil
	}

	// 获得不同驱动的存储
	var stored store.StoreInterface
	switch Config.Data.Cache.Engine {
	case "redis":
		cnf := Config.Data.Cache.Redis
		if len(cnf.Addr) == 0 || !strings.Contains(cnf.Addr, ":") {
			return xerror.New("cache redis config not exist")
		}

		stored = redisstore.NewRedis(
			redis.NewClient(&redis.Options{
				Addr:     cnf.Addr,
				Password: cnf.Password,
				DB:       cnf.DB,
			}),
			store.WithExpiration(15*time.Minute), // 默认缓存15分钟，防止缓存默认永存
		)
	default:
		return xerror.New("cache drive not support")
	}

	ctx := context.Background()
	cacheManager := cache.New[string](stored)
	// 设置测试缓存
	if err := cacheManager.Set(ctx, "cacheMangerTest", "test cache",
		// 显示指定过期时间，覆盖 store 的默认行为
		store.WithExpiration(10*time.Minute)); nil != err {
		return xerror.Wrap(err, "cache manager failed")
	}

	_cacheManager = cacheManager
	Log.Debug("cache manger init ... success")
	return nil
}

func initLocker() error {
	if !Config.Data.Locker.Enabled {
		Log.Debug("locker disabled, skip")
		return nil
	}

	var (
		err  error
		lock locker.ILocker
	)
	switch Config.Data.Locker.Engine {
	case "redis":
		cnf := Config.Data.Locker.Redis
		lock, err = inlocker.RedisLocker(&redis.Options{
			Addr:     cnf.Addr,
			Password: cnf.Password,
			DB:       cnf.DB,
		})
	default:
		return xerror.New("locker drive not support")
	}
	if nil != err {
		return err
	}

	_locker = lock
	Log.Debug("locker enabled, init success")

	return nil
}

func MigrateData() error {
	// 初始化数据
	if err := migrateData(); nil != err {
		return xerror.Wrap(err, "data builder init failed")
	}

	return nil
}

func migrateData() error {
	if !Config.Data.Persistent.Enabled || !Config.Data.Persistent.AutoMigrate {
		Log.Debug("database auto migrate disabled, skip")
		return nil
	}

	if err := database.Migrate(Database(), Config); nil != err {
		return err
	}

	return nil
}
