package redis

import (
	"19shubham11/url-shortener/cmd/conf"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func SetupRedis(config conf.RedisConf) *redis.Client {
	// redisUrl := fmt.Sprintf("redis://%s:%s@%s:%d/%d", config.Username, config.Password, config.Host, config.Port, config.DB)
	_ = fmt.Sprintf("%s:%d", config.Host, config.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	})
	return rdb
}
