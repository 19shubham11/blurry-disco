package redis

import (
	redisClient "github.com/gomodule/redigo/redis"
)

type RedisModel struct {
	Redis redisClient.Conn
}

func (r RedisModel) Set(key string, value string) (string, error) {
	ok, err := redisClient.String(r.Redis.Do("SET", key, value))
	if err != nil {
		return "", err
	}
	return ok, nil
}

func (r RedisModel) Get(key string) (string, error) {
	value, err := redisClient.String(r.Redis.Do("GET", key))
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r RedisModel) Incr(key string) (int, error) {
	value, err := redisClient.Int(r.Redis.Do("INCR", key))
	if err != nil {
		return -1, err
	}
	return value, nil
}
