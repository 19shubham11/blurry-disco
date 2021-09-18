package redis

import (
	"19shubham11/url-shortener/cmd/config"

	"github.com/go-redis/redis/v8"
)

func SetupRedis(config config.RedisConf) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	})

	return rdb
}
