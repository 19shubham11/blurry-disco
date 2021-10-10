package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Model struct {
	Redis *redis.Client
	Ctx   context.Context
}

func (r Model) Set(key string, value string) error {
	err := r.Redis.Set(r.Ctx, key, value, 0).Err()
	if err != nil {
		return ErrorInternal
	}

	return nil
}

func (r Model) Get(key string) (string, error) {
	value, err := r.Redis.Get(r.Ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return value, ErrorNotFound
		}

		return "", ErrorInternal
	}

	return value, nil
}

func (r Model) Incr(key string) (int, error) {
	value, err := r.Redis.Incr(r.Ctx, key).Result()
	if err != nil {
		return 0, ErrorInternal
	}

	return int(value), nil
}

func (r Model) Mset(set map[string]string) error {
	redisSet := []string{}

	for key, val := range set {
		redisSet = append(redisSet, key, val)
	}

	err := r.Redis.MSet(r.Ctx, redisSet).Err()
	if err != nil {
		return ErrorInternal
	}

	return nil
}

func (r Model) Mget(keys []string) ([]string, error) {
	res, err := r.Redis.MGet(r.Ctx, keys...).Result()

	if err != nil {
		return nil, ErrorInternal
	}

	ret := []string{}
	for _, v := range res {
		ret = append(ret, fmt.Sprint(v))
	}

	return ret, nil
}
