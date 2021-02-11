package redis

import (
	config "19shubham11/url-shortener/cmd/conf"
	"log"
	"os"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

var conn redis.Conn
var redisModel *RedisModel

func redisSetup() (redis.Conn, func()) {

	redisPass := os.Getenv("REDIS_PASS")

	redisConf := config.RedisConf{
		Host:     "localhost",
		Port:     6379,
		Username: "default",
		Password: redisPass,
		DB:       3,
	}

	conn, err := SetupRedis(redisConf)
	if err != nil {
		log.Fatalf("Redis Error!")
	}

	return conn, func() {
		conn.Do("FLUSHDB")
		conn.Close()
	}
}

func TestMain(m *testing.M) {
	conn, teardown := redisSetup()
	defer teardown()

	redisModel = &RedisModel{Redis: conn}
	code := m.Run()
	os.Exit(code)
}

func TestSET(t *testing.T) {
	t.Run("SET should return OK", func(t *testing.T) {
		str, err := redisModel.Set("key1", "value1")

		assert.Nil(t, err)
		assert.Equal(t, "OK", str)
	})
}
func TestGET(t *testing.T) {
	t.Run("GET should return the value of a key if it exists", func(t *testing.T) {
		key := "key1"
		value := "value1"
		_, err := redisModel.Set(key, value)
		if err != nil {
			t.Fatal(err)
		}

		res, err := redisModel.Get(key)
		assert.Nil(t, err)
		assert.Equal(t, res, value)
	})

	t.Run("GET should return an error if trying to get a key that does not exist", func(t *testing.T) {
		res, err := redisModel.Get("Does not exist")
		assert.Equal(t, res, "")
		assert.Equal(t, err.Error(), "redigo: nil returned")
	})
}
func TestINCR(t *testing.T) {
	t.Run("INCR should return the incremented value if the value can be cast as int", func(t *testing.T) {
		key := "key1"
		value := "22"
		_, err := redisModel.Set(key, value)
		if err != nil {
			t.Fatal(err)
		}

		res, err := redisModel.Incr(key)
		assert.Nil(t, err)
		assert.Equal(t, res, 23)
	})

	t.Run("INCR should return an error if the value cannot be cast as int", func(t *testing.T) {
		key := "key1"
		value := "value"
		_, err := redisModel.Set(key, value)
		if err != nil {
			t.Fatal(err)
		}

		_, err = redisModel.Incr(key)
		assert.Equal(t, err.Error(), "ERR value is not an integer or out of range")
	})
}
