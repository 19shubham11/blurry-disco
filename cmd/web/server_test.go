package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	config "19shubham11/url-shortener/cmd/conf"
	db "19shubham11/url-shortener/pkg/redis"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func redisSetup() (redis.Conn, func()) {
	redisPass := os.Getenv("REDIS_PASS")

	redisConf := config.RedisConf{
		Host:     "localhost",
		Port:     6379,
		Username: "default",
		Password: redisPass,
		DB:       3,
	}

	conn, err := db.SetupRedis(redisConf)
	if err != nil {
		log.Fatal(err)
	}

	return conn, func() {
		conn.Do("FLUSHDB")
		conn.Close()
	}
}

func getResponseTextBody(resp *http.Response) []byte {
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData
}

func getJSONBytes(str interface{}) []byte {
	b, err := json.Marshal(str)

	if err != nil {
		log.Fatal(err)
	}
	return b
}

var ts *httptest.Server
var app *application
var client *http.Client

func TestMain(m *testing.M) {
	conn, teardown := redisSetup()
	defer teardown()
	log.Println("Tests - Connected to Redis!")

	redisModel := db.RedisModel{Redis: conn}
	app = &application{DB: redisModel, BaseURL: "myURLShortener"}
	app.routes()

	ts = httptest.NewServer(app.routes())

	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	defer ts.Close()
	code := m.Run()
	os.Exit(code)
}

func TestGETHealth(t *testing.T) {
	t.Run("Should return 200 OK", func(t *testing.T) {
		resp, err := client.Get(ts.URL + "/internal/health")
		if err != nil {
			log.Fatal(err)
		}

		respBody := getResponseTextBody(resp)
		respString := string(respBody)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, "OK", respString)
	})
}

func TestPOSTShorten(t *testing.T) {
	t.Run("Should return 200 for a valid request", func(t *testing.T) {
		reqBody := &ShortenURLRequest{URL: "http://www.google.com"}
		reqJSON := getJSONBytes(reqBody)
		resp, err := client.Post(ts.URL+"/shorten", "application/json", bytes.NewBuffer(reqJSON))

		if err != nil {
			log.Fatal(err)
		}

		response := &ShortenURLResponse{}
		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(response)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, true, strings.Contains(response.ShortenedURL, "myURLShortener/"))

	})

	t.Run("Should return 400 if the given url is invalid", func(t *testing.T) {
		reqBody := &ShortenURLRequest{URL: "notAURL"}
		reqJSON := getJSONBytes(reqBody)
		resp, err := client.Post(ts.URL+"/shorten", "application/json", bytes.NewBuffer(reqJSON))

		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("Should return 400 if the input does not contain the key `url`", func(t *testing.T) {
		reqBody := []byte(`{"key: "http://www.google.com"}`)
		resp, err := client.Post(ts.URL+"/shorten", "application/json", bytes.NewBuffer(reqBody))

		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 400, resp.StatusCode)
	})
}

func TestGETOriginal(t *testing.T) {
	t.Run("Should redirect to the original URL if the given ID is valid", func(t *testing.T) {
		// set a key/value pair in redis
		key := "205db389"
		app.DB.Set(key, "https://www.google.com")
		reqURL := fmt.Sprintf("%s/%s", ts.URL, key)
		resp, err := client.Get(reqURL)

		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 302, resp.StatusCode)
	})

	t.Run("Should return 404 if the given ID does not exist", func(t *testing.T) {
		reqURL := fmt.Sprintf(ts.URL + "/21345")
		resp, err := client.Get(reqURL)

		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, 404, resp.StatusCode)
	})
}
