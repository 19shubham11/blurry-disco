package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-redis/redis/v8"

	"19shubham11/url-shortener/cmd/config"
	"19shubham11/url-shortener/cmd/pkg/redis"
	"19shubham11/url-shortener/cmd/pkg/store"
)

type application struct {
	DB      store.Store
	BaseURL string
}

func main() {
	appConfig := config.GetApplicationConfig()
	conn := redis.Setup(appConfig.Redis)

	_, err := conn.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Unable to connect to redis!", err)
	}

	ctx := context.Background()

	log.Printf("Connected to Redis on %s:%d", appConfig.Redis.Host, appConfig.Server.Port)

	baseURL := fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port)

	redisModel := redis.Model{Redis: conn, Ctx: ctx}
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
