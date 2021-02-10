package redis

import (
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func RedisSetup(t *testing.T) (redis.Conn, func()) {
	t.Helper()
	conn, err := SetupRedis()
	if err != nil {
		t.Fatalf("Redis Error!")
	}

	return conn, func() {
		conn.Do("FLUSHDB")
		conn.Close()
	}
}

func TestRedis(t *testing.T) {
	conn, teardown := RedisSetup(t)
	defer teardown()

	redisModel := &RedisModel{Redis: conn}

	t.Run("SET should return OK", func(t *testing.T) {
		str, err := redisModel.Set("key1", "value1")

		assert.Nil(t, err)
		assert.Equal(t, "OK", str)
	})

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
