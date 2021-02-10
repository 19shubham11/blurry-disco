package redis

import "github.com/gomodule/redigo/redis"

func SetupRedis() (redis.Conn, error) {
	redisUrl := "redis://localhost:6379/2"

	conn, err := redis.DialURL(redisUrl)

	return conn, err
}
