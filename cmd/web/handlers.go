package main

import (
	helpers "19shubham11/url-shortener/cmd/helpers"
	"encoding/json"
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

	shortenedURLResponse := app.shortenURLController(body)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shortenedURLResponse)
}

func (app *application) getOriginalURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["id"]

	url, err := app.getOriginalURLController(hash)

	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
