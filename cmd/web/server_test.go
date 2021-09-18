package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"19shubham11/url-shortener/cmd/config"
	db "19shubham11/url-shortener/cmd/pkg/redis"
)

func redisSetup() (*redis.Client, func()) {
	redisPass := os.Getenv("REDIS_PASS")

	redisConf := config.RedisConf{
		Host:     "localhost",
		Port:     6379,
		Username: "default",
		Password: redisPass,
		DB:       3,
	}

	conn := db.SetupRedis(redisConf)

	return conn, func() {
		var ctx = context.Background()

		conn.FlushDB(ctx)
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

func initTests(m *testing.M) int {
	conn, teardown := redisSetup()
	defer teardown()

	log.Println("Tests - Connected to Redis!")

	ctx := context.Background()

	redisModel := db.Model{Redis: conn, Ctx: ctx}
	app = &application{DB: redisModel}

	ts = httptest.NewServer(app.routes())
	app.BaseURL = ts.URL

	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	defer ts.Close()

	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(initTests(m))
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
		reqBody := &ShortenURLRequest{URL: "https://www.google.com"}
		reqJSON := getJSONBytes(reqBody)
		resp, err := client.Post(ts.URL+"/shorten", "application/json", bytes.NewBuffer(reqJSON))

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		response := &ShortenURLResponse{}
		decoder := json.NewDecoder(resp.Body)
		_ = decoder.Decode(response)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, true, strings.Contains(response.ShortenedURL, ts.URL))
	})

	t.Run("Should return 400 if url field is empty", func(t *testing.T) {
		reqBody := &ShortenURLRequest{URL: ""}
		reqJSON := getJSONBytes(reqBody)
		resp, err := client.Post(ts.URL+"/shorten", "application/json", bytes.NewBuffer(reqJSON))

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Should return 400 if the given url is invalid", func(t *testing.T) {
		reqBody := &ShortenURLRequest{URL: "notAURL"}
		reqJSON := getJSONBytes(reqBody)
		resp, err := client.Post(ts.URL+"/shorten", "application/json", bytes.NewBuffer(reqJSON))

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Should return 400 if the input does not contain the key `url`", func(t *testing.T) {
		reqBody := []byte(`{"key: "http://www.google.com"}`)
		resp, err := client.Post(ts.URL+"/shorten", "application/json", bytes.NewBuffer(reqBody))

		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestGETOriginal(t *testing.T) {
	t.Run("Should redirect to the original URL if the given ID is valid", func(t *testing.T) {
		// set a key/value pair in redis
		key := "205db389"
		_ = app.DB.Set(key, "https://www.google.com")
		reqURL := fmt.Sprintf("%s/%s", ts.URL, key)
		resp, err := client.Get(reqURL)

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		assert.Equal(t, http.StatusFound, resp.StatusCode)
	})

	t.Run("Should return 404 if the given ID does not exist", func(t *testing.T) {
		reqURL := fmt.Sprintf(ts.URL + "/21345")
		resp, err := client.Get(reqURL)

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestGetStats(t *testing.T) {
	t.Run("Should return 200 and stats for a given id", func(t *testing.T) {
		url := "https://www.google.com"
		reqBody := &ShortenURLRequest{URL: url}
		reqJSON := getJSONBytes(reqBody)
		resp, err := client.Post(ts.URL+"/shorten", "application/json", bytes.NewBuffer(reqJSON))

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		shortenURLResponse := &ShortenURLResponse{}
		decoder := json.NewDecoder(resp.Body)
		_ = decoder.Decode(shortenURLResponse)

		numOfRequests := 50
		makeConcurrentHTTPCalls(shortenURLResponse.ShortenedURL, numOfRequests)

		hitsRes, err := client.Get(shortenURLResponse.ShortenedURL + "/stats")
		if err != nil {
			log.Fatal(err)
		}

		defer hitsRes.Body.Close()

		statsResponse := &StatsResponse{}
		expectedResp := &StatsResponse{URL: url, Hits: numOfRequests}

		decoder = json.NewDecoder(hitsRes.Body)
		_ = decoder.Decode(statsResponse)

		assert.Equal(t, http.StatusOK, hitsRes.StatusCode)
		assert.Equal(t, expectedResp, statsResponse)
	})

	t.Run("Should return 404 for an id that does not exist", func(t *testing.T) {
		reqURL := fmt.Sprintf(ts.URL + "/21345s/stats")
		resp, err := client.Get(reqURL)

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func makeConcurrentHTTPCalls(url string, noOfRequests int) {
	var wg sync.WaitGroup

	for i := 0; i < noOfRequests; i++ {
		wg.Add(1)

		makeClientReq := func() {
			defer wg.Done()

			resp, err := client.Get(url)
			if err != nil {
				log.Fatal(err)
			}

			defer resp.Body.Close()
		}

		go makeClientReq()
	}

	wg.Wait()
}
