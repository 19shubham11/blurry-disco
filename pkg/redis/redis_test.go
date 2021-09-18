package redis

import (
	"context"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"19shubham11/url-shortener/config"
	customerrors "19shubham11/url-shortener/internal/customerrors"
)

var conn *redis.Client
var redisModel *Model

func redisSetup() (*redis.Client, func()) {
	redisPass := os.Getenv("REDIS_PASS")

	redisConf := config.RedisConf{
		Host:     "localhost",
		Port:     6379,
		Username: "default",
		Password: redisPass,
		DB:       3,
	}

	conn = SetupRedis(redisConf)

	return conn, func() {
		var ctx = context.Background()

		conn.FlushDB(ctx)
		conn.Close()
	}
}

func initTests(m *testing.M) int {
	conn, teardown := redisSetup()
	defer teardown()

	ctx := context.Background()
	redisModel = &Model{Redis: conn, Ctx: ctx}

	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(initTests(m))
}

func TestSET(t *testing.T) {
	t.Run("SET should return OK", func(t *testing.T) {
		err := redisModel.Set("key1", "value1")
		assert.Equal(t, err, nil)
	})
}
func TestGET(t *testing.T) {
	t.Run("GET should return the value of a key if it exists", func(t *testing.T) {
		key := "keyGet"
		value := "value1"

		err := redisModel.Set(key, value)
		assert.Nil(t, err)

		res, err := redisModel.Get(key)
		assert.Nil(t, err)
		assert.Equal(t, res, value)
	})

	t.Run("GET should return an error if trying to get a key that does not exist", func(t *testing.T) {
		res, err := redisModel.Get("Does not exist")
		assert.Equal(t, res, "")
		assert.Equal(t, err, customerrors.ErrorNotFound)
	})
}
func TestINCR(t *testing.T) {
	t.Run("INCR should return the incremented value if the value can be cast as int", func(t *testing.T) {
		key := "keyIncr"
		value := "22"

		err := redisModel.Set(key, value)
		assert.Nil(t, err)

		res, err := redisModel.Incr(key)
		assert.Nil(t, err)
		assert.Equal(t, res, 23)
	})

	t.Run("INCR should return an error if the value cannot be cast as int", func(t *testing.T) {
		key := "keyNotInt"
		value := "value"

		err := redisModel.Set(key, value)
		assert.Nil(t, err)

		_, err = redisModel.Incr(key)
		assert.Equal(t, err, customerrors.ErrorInternal)
	})
}
