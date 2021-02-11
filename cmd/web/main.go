package main

import (
	"19shubham11/url-shortener/pkg/redis"
	"log"
	"net/http"
)

type application struct {
	DB      redis.Store
	BaseURL string
}

func main() {
	conn, err := redis.SetupRedis()
	if err != nil {
		log.Fatalf("Redis Error! %v", err)
	}

	redisModel := redis.RedisModel{Redis: conn}
	app := &application{DB: redisModel, BaseURL: "localhost:2001"}

	server := &http.Server{
		Addr:    ":2001",
		Handler: app.routes(),
	}

	log.Println("Starting server!")
	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("Something Happened!")
	}
}
