package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/internal/health", app.checkHealth).Methods("GET")

	router.HandleFunc("/shorten", app.shortenURL).Methods("POST")
	router.HandleFunc("/{id}", app.getOriginalURL).Methods("GET")
	router.HandleFunc("/{id}/stats", app.getStats).Methods("GET")

	return router
}
