package redis

import (
	"github.com/gomodule/redigo/redis"
)

type RedisModel struct {
	Redis redis.Conn
}

// Set saves a key value pair in redis, returns OK or an error
func (r RedisModel) Set(key string, value string) (string, error) {
	ok, err := redis.String(r.Redis.Do("SET", key, value))
	if err != nil {
		return "", err
	}
	return ok, nil
}

// Get returns the value for a given key, or an error if the key does not exist
func (r RedisModel) Get(key string) (string, error) {
	value, err := redis.String(r.Redis.Do("GET", key))
	if err != nil {
		return "", err
	}
	return value, nil
}
