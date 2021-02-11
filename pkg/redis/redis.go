package redis

import (
	"github.com/gomodule/redigo/redis"
)

type RedisModel struct {
	Redis redis.Conn
}

func (r RedisModel) Set(key string, value string) (string, error) {
	ok, err := redis.String(r.Redis.Do("SET", key, value))
	if err != nil {
		return "", err
	}
	return ok, nil
}

func (r RedisModel) Get(key string) (string, error) {
	value, err := redis.String(r.Redis.Do("GET", key))
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r RedisModel) Incr(key string) (int, error) {
	value, err := redis.Int(r.Redis.Do("INCR", key))
	if err != nil {
		return -1, err
	}
	return value, nil
}
