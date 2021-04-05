package main

import (
	customErrors "19shubham11/url-shortener/cmd/customErrors"
	helpers "19shubham11/url-shortener/cmd/helpers"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) checkHealth(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) shortenURL(w http.ResponseWriter, r *http.Request) {
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

	res, err := app.shortenURLController(body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (app *application) getOriginalURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["id"]

	url, err := app.getOriginalURLController(hash)

	if err != nil {
		if errors.Is(err, customErrors.ErrorNotFound) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		} else if errors.Is(err, customErrors.ErrorInternal) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (app *application) getStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["id"]

	res, err := app.getStatsController(hash)

	if err != nil {
		if errors.Is(err, customErrors.ErrorNotFound) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		} else if errors.Is(err, customErrors.ErrorInternal) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
