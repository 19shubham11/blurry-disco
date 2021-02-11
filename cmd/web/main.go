package main

import (
	"19shubham11/url-shortener/pkg/redis"
	"19shubham11/url-shortener/pkg/store"
	"log"
	"net/http"
)

type application struct {
	DB      store.Store
	BaseURL string
}

func main() {
	conn, err := redis.SetupRedis()
	if err != nil {
		log.Fatalf("Redis Connection Error! %v", err)
	}
	log.Println("Connected to Redis!")

	redisModel := redis.RedisModel{Redis: conn}
	app := &application{DB: redisModel, BaseURL: "localhost:2001"}

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
