package redis

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"

	"19shubham11/url-shortener/cmd/internal/customerrors"
)

type Model struct {
	Redis *redis.Client
	Ctx   context.Context
}

func (r Model) Set(key string, value string) error {
	err := r.Redis.Set(r.Ctx, key, value, 0).Err()
	if err != nil {
		return customerrors.ErrorInternal
	}

	return nil
}

func (r Model) Get(key string) (string, error) {
	value, err := r.Redis.Get(r.Ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return value, customerrors.ErrorNotFound
		}

		return "", customerrors.ErrorInternal
	}

	return value, nil
}

func (r Model) Incr(key string) (int, error) {
	value, err := r.Redis.Incr(r.Ctx, key).Result()
	if err != nil {
		return 0, customerrors.ErrorInternal
	}

	return int(value), nil
}
