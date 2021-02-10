package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	routes := mux.NewRouter()
	routes.HandleFunc("/internal/health", app.health).Methods("GET")

	return routes
}
