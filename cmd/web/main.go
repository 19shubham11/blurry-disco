package main

import (
	config "19shubham11/url-shortener/config"
	"19shubham11/url-shortener/pkg/redis"
	"19shubham11/url-shortener/pkg/store"
	"context"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-redis/redis/v8"
)

type application struct {
	DB      store.Store
	BaseURL string
}

func main() {
	appConfig := config.GetApplicationConfig()
	conn := redis.SetupRedis(appConfig.Redis)

	_, err := conn.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Unable to connect to redis!", err)
	}

	ctx := context.Background()
	log.Println("Connected to Redis!")

	redisModel := redis.RedisModel{Redis: conn, Ctx: ctx}
	baseURL := fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port)
	app := &application{DB: redisModel, BaseURL: baseURL}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", appConfig.Server.Port),
		Handler: app.routes(),
	}

	log.Printf("Starting server on port %d!", appConfig.Server.Port)
	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("Something Happened!")
	}
}
