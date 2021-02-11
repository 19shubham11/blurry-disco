package redis

import (
	"19shubham11/url-shortener/cmd/conf"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func SetupRedis(config conf.RedisConf) (redis.Conn, error) {
	redisUrl := fmt.Sprintf("redis://%s:%s@%s:%d/%d", config.Username, config.Password, config.Host, config.Port, config.DB)

	return redis.DialURL(redisUrl)
}
