package redis

import (
	"19shubham11/url-shortener/cmd/conf"
	"fmt"

	redisClient "github.com/gomodule/redigo/redis"
)

func SetupRedis(config conf.RedisConf) (redisClient.Conn, error) {
	redisURL := fmt.Sprintf("redis://%s:%s@%s:%d/%d", config.Username, config.Password, config.Host, config.Port, config.DB)

	return redisClient.DialURL(redisURL)
}
