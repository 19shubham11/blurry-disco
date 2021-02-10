package main

import (
	"log"
	"net/http"
)

type application struct{}

func main() {
	app := &application{}

	server := &http.Server{
		Addr:    ":2001",
		Handler: app.routes(),
	}

	log.Println("Starting server!")
	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("Something Happened!")
	}
}
