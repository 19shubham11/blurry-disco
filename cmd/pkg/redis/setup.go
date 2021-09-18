package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"

	"19shubham11/url-shortener/cmd/config"
)

func Setup(config config.RedisConf) *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	})

	return redis
}
