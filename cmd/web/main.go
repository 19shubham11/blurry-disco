package main

import (
	config "19shubham11/url-shortener/cmd/conf"
	"19shubham11/url-shortener/pkg/redis"
	"19shubham11/url-shortener/pkg/store"
	"fmt"
	"log"
	"net/http"
)

type application struct {
	DB      store.Store
	BaseURL string
}

func main() {
	appConfig := config.GetApplicationConfig()
	conn, err := redis.SetupRedis(appConfig.Redis)
	if err != nil {
		log.Fatalf("Redis Connection Error! %v", err)
	}
	log.Println("Connected to Redis!")

	redisModel := redis.RedisModel{Redis: conn}
	baseURL := fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port)
	app := &application{DB: redisModel, BaseURL: baseURL}

	server := &http.Server{
		Addr:    ":2001",
		Handler: app.routes(),
	}

	log.Println("Starting server on port 2001!")
	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("Something Happened!")
	}
}
