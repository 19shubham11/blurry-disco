package redis

import (
	"context"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"19shubham11/url-shortener/cmd/config"
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

	conn = Setup(redisConf)

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
	err := redisModel.Set("key1", "value1")
	assert.Equal(t, err, nil)
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
		assert.Equal(t, err, ErrorNotFound)
	})
}

func TestINCR(t *testing.T) {
	tests := map[string]struct {
		key           string
		value         string
		expectedError error
		expectedRes   int
	}{
		"value cast as int":           {"KeyIncr", "22", nil, 23},
		"value cannot be cast as int": {"keyNotInt", "twentyTwo", ErrorInternal, 0},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := redisModel.Set(tc.key, tc.value)
			assert.Nil(t, err)

			res, err := redisModel.Incr(tc.key)

			if err != nil {
				assert.Equal(t, err, ErrorInternal)
			} else {
				assert.Equal(t, res, tc.expectedRes)
			}
		})
	}
}

func TestMset(t *testing.T) {
	set := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}

	err := redisModel.Mset(set)
	assert.Equal(t, err, nil)
}

func TestMget(t *testing.T) {
	set := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}

	values := make([]string, 0, len(set))
	keys := make([]string, 0, len(set))

	for k, v := range set {
		keys = append(keys, k)
		values = append(values, v)
	}

	_ = redisModel.Mset(set)

	res, err := redisModel.Mget(keys)

	assert.Nil(t, err)
	assert.Equal(t, res, values)
}
