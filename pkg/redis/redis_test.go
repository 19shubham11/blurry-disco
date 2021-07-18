package redis

import (
	config "19shubham11/url-shortener/config"
	"19shubham11/url-shortener/internal/customErrors"
	"context"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var conn *redis.Client
var redisModel *RedisModel

func redisSetup() (*redis.Client, func()) {

	redisPass := os.Getenv("REDIS_PASS")

	redisConf := config.RedisConf{
		Host:     "localhost",
		Port:     6379,
		Username: "default",
		Password: redisPass,
		DB:       3,
	}

	conn := SetupRedis(redisConf)

	return conn, func() {
		var ctx = context.Background()
		conn.FlushDB(ctx)
		conn.Close()
	}
}

func TestMain(m *testing.M) {
	conn, teardown := redisSetup()
	defer teardown()
	ctx := context.Background()
	redisModel = &RedisModel{Redis: conn, Ctx: ctx}
	code := m.Run()
	os.Exit(code)
}

func TestSET(t *testing.T) {
	t.Run("SET should return OK", func(t *testing.T) {
		err := redisModel.Set("key1", "value1")
		assert.Equal(t, err, nil)
	})
}
func TestGET(t *testing.T) {
	t.Run("GET should return the value of a key if it exists", func(t *testing.T) {
		key := "key1"
		value := "value1"
		redisModel.Set(key, value)

		res, err := redisModel.Get(key)
		assert.Nil(t, err)
		assert.Equal(t, res, value)
	})

	t.Run("GET should return an error if trying to get a key that does not exist", func(t *testing.T) {
		res, err := redisModel.Get("Does not exist")
		assert.Equal(t, res, "")
		assert.Equal(t, err, customErrors.ErrorNotFound)
	})
}
func TestINCR(t *testing.T) {
	t.Run("INCR should return the incremented value if the value can be cast as int", func(t *testing.T) {
		key := "key1"
		value := "22"
		redisModel.Set(key, value)

		res, err := redisModel.Incr(key)
		assert.Nil(t, err)
		assert.Equal(t, res, 23)
	})

	t.Run("INCR should return an error if the value cannot be cast as int", func(t *testing.T) {
		key := "key1"
		value := "value"
		redisModel.Set(key, value)

		_, err := redisModel.Incr(key)
		assert.Equal(t, err, customErrors.ErrorInternal)
	})
}
