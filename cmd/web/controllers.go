package main

import (
	"fmt"
	"strconv"

	"19shubham11/url-shortener/cmd/internal/helpers"
)

const (
	statsPrefix = "STATS"
	initStats   = 0
)

func (app *application) shortenURLController(request *ShortenURLRequest) (ShortenURLResponse, error) {
	hash := helpers.CreateUniqueHash()
	statKey := getStatsKey(hash)

	err := app.DB.Set(hash, request.URL)
	if err != nil {
		return ShortenURLResponse{}, err
	}

	err = app.DB.Set(statKey, strconv.Itoa(initStats))

	if err != nil {
		return ShortenURLResponse{}, err
	}

	returnURL := fmt.Sprintf("%s/%s", app.BaseURL, hash)

	return ShortenURLResponse{ShortenedURL: returnURL}, nil
}

func (app *application) getOriginalURLController(hash string) (url string, err error) {
	url, err = app.DB.Get(hash)
	if err != nil {
		return "", err
	}

	statKey := getStatsKey(hash)
	_, err = app.DB.Incr(statKey)

	if err != nil {
		return "", err
	}

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

	hitsInt, err := strconv.Atoi(hits)

	if err != nil {
		return StatsResponse{}, err
	}

	res = StatsResponse{URL: url, Hits: hitsInt}

	return res, nil
}

func getStatsKey(hash string) string {
	return fmt.Sprintf("%s_%s", statsPrefix, hash)
}
