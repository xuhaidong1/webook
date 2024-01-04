package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/xuhaidong1/webook/internal/domain"
	"time"
)

// ErrKeyNotExist 因为我们目前还是只有一个实现，所以可以保持用别名
var ErrKeyNotExist = redis.Nil

//go:generate mockgen -source=./user.go -package=cachemocks -destination=mocks/user.mock.go UserCache
type UserCache interface {
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type RedisUserCache struct {
	cmd redis.Cmdable
	// 过期时间
	expiration time.Duration
}

func NewRedisUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
func (cache *RedisUserCache) Delete(ctx context.Context, id int64) error {
	return cache.cmd.Del(ctx, cache.key(id)).Err()
}

func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	data, err := cache.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	// 反序列化回来
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	return cache.cmd.Set(ctx, key, data, cache.expiration).Err()
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
