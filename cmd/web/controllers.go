package main

import (
	"19shubham11/url-shortener/cmd/helpers"
	"fmt"
	"strconv"
)

const statsPrefix = "STATS"
const initStats = 0

func (app *application) shortenURLController(request *ShortenURLRequest) ShortenURLResponse {
	hash := helpers.CreateRandomHash()
	statKey := getStatsKey(hash)

	app.DB.Set(hash, request.URL)
	app.DB.Set(statKey, strconv.Itoa(initStats))

	returnURL := fmt.Sprintf("%s/%s", app.BaseURL, hash)

	return ShortenURLResponse{ShortenedURL: returnURL}
}

func (app *application) getOriginalURLController(hash string) (url string, err error) {
	url, err = app.DB.Get(hash)
	statKey := getStatsKey(hash)
	if err != nil {
		return "", err
	}
	app.DB.Incr(statKey)
	return url, err
}

func (app *application) getStatsController(hash string) (res StatsResponse, err error) {
	url, err := app.DB.Get(hash)
	if err != nil {
		return StatsResponse{}, err
	}

	statKey := getStatsKey(hash)
	hits, err := app.DB.Get(statKey)
	if err != nil {
		return StatsResponse{}, err
	}

	res = StatsResponse{URL: url, Hits: hits}
	return res, nil
}

func getStatsKey(hash string) string {
	return fmt.Sprintf("%s_%s", statsPrefix, hash)
}
