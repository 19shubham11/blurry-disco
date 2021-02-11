package main

import (
	"19shubham11/url-shortener/cmd/helpers"
	"fmt"
)

const statsPrefix = "STATS"
const initStats = 0

func (app *application) shortenURLController(request *ShortenURLRequest) ShortenURLResponse {
	hash := helpers.CreateRandomHash()
	// statKey := getStatsKey(hash)

	app.DB.Set(hash, request.URL)
	// app.DB.Set(statKey, string(initStats))

	returnURL := fmt.Sprintf("%s/%s", app.BaseURL, hash)

	return ShortenURLResponse{ShortenedURL: returnURL}
}

func (app *application) getOriginalURLController(hash string) (url string, err error) {
	return app.DB.Get(hash)
}

func getStatsKey(hash string) string {
	return fmt.Sprintf("%s_%s", statsPrefix, hash)
}
