package helper

import (
	"server-api/global"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"

	"github.com/gomooth/pkg/http/jwt"
	"github.com/gomooth/pkg/http/middleware"
)

// JWTOption JWT 相关配置
func JWTOption(refresh bool) *jwt.Option {
	opt := &jwt.Option{
		RoleConvert:     global.NewRole,
		RefreshDuration: 0, // 0-不自动刷新
		Secret:          []byte(global.Config.App.Secret),
	}

	refreshDuration := time.Duration(0)
	if refresh {
		refreshDuration = 12 * time.Hour
	}

	opt.RefreshDuration = refreshDuration

	return opt
}

// SessionStore session 存储
func SessionStore(opt middleware.SessionOption) sessions.Store {
	if global.Config.Data.Redis.Enabled {
		return sessionRedisStore(opt)
	}

	return memstore.NewStore([]byte(global.Config.App.Secret))
}

func sessionRedisStore(opt middleware.SessionOption) sessions.Store {
	store, err := redis.NewStore(
		int(opt.MaxAge.Minutes()), // 有效时间，分钟
		"tcp",
		global.Config.Data.Redis.Addr,
		global.Config.Data.Redis.Password,
		[]byte(global.Config.App.Secret),
	)
	if nil != err {
		global.Log.Errorf("session redis store failed: %+v", err)
	}

	return store
}
