package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisModel struct {
	Redis *redis.Client
}

func (r RedisModel) Set(key string, value string) error {
	return r.Redis.Set(ctx, key, value, 0).Err()

}

func (r RedisModel) Get(key string) (string, error) {
	value, err := r.Redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r RedisModel) Incr(key string) (int, error) {
	value, err := r.Redis.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(value), nil
}
