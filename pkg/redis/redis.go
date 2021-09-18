package redis

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"

	"19shubham11/url-shortener/internal/customErrors"
)

type RedisModel struct {
	Redis *redis.Client
	Ctx   context.Context
}

func (r RedisModel) Set(key string, value string) error {
	err := r.Redis.Set(r.Ctx, key, value, 0).Err()
	if err != nil {
		return customErrors.ErrorInternal
	}
	return nil
}

func (r RedisModel) Get(key string) (string, error) {
	value, err := r.Redis.Get(r.Ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return value, customErrors.ErrorNotFound
		} else {
			return "", customErrors.ErrorInternal
		}
	}
	return value, nil
}

func (r RedisModel) Incr(key string) (int, error) {
	value, err := r.Redis.Incr(r.Ctx, key).Result()
	if err != nil {
		return 0, customErrors.ErrorInternal
	}
	return int(value), nil
}
