package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"19shubham11/url-shortener/cmd/internal/helpers"
	"19shubham11/url-shortener/cmd/pkg/redis"
)

func (s *Server) checkHealth(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("OK"))
}

func (s *Server) shortenURL(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	body := &ShortenURLRequest{}

	err := decoder.Decode(body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if body.URL == "" {
		http.Error(w, "Missing required field `URL`", http.StatusBadRequest)
		return
	}

	if !helpers.IsValidURL(body.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	res, err := s.shortenURLController(body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) getOriginalURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["id"]

	url, err := s.getOriginalURLController(hash)

	if err != nil {
		if errors.Is(err, redis.ErrorNotFound) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		} else if errors.Is(err, redis.ErrorInternal) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (s *Server) getStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["id"]

	res, err := s.getStatsController(hash)

	if err != nil {
		if errors.Is(err, redis.ErrorNotFound) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		} else if errors.Is(err, redis.ErrorInternal) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
