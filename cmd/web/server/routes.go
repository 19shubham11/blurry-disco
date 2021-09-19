package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) Routes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/internal/health", s.checkHealth).Methods("GET")
	router.HandleFunc("/shorten", s.shortenURL).Methods("POST")
	router.HandleFunc("/{id}", s.getOriginalURL).Methods("GET")
	router.HandleFunc("/{id}/stats", s.getStats).Methods("GET")

	return router
}
