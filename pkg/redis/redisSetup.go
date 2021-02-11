package redis

import "github.com/gomodule/redigo/redis"

func SetupRedis() (redis.Conn, error) {
	redisUrl := "redis://localhost:6379/2"

	return redis.DialURL(redisUrl)
}
