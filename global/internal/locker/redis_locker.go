package locker

import (
	"context"
	"strings"

	"github.com/gomooth/locker"
	"github.com/gomooth/locker/redislock"

	"github.com/redis/go-redis/v9"

	"github.com/gomooth/xerror"
)

func RedisLocker(opt *redis.Options) (locker.ILocker, *redis.Client, error) {
	if len(opt.Addr) == 0 || !strings.Contains(opt.Addr, ":") {
		return nil, nil, xerror.New("locker redis config not exist")
	}

	client := redis.NewClient(opt)
	if err := client.Ping(context.Background()).Err(); nil != err {
		_ = client.Close()
		return nil, nil, xerror.Wrap(err, "redis client connect failed")
	}

	return redislock.New(client), client, nil
}
