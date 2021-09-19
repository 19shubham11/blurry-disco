package server

import (
	"fmt"
	"strconv"

	"19shubham11/url-shortener/cmd/internal/helpers"
)

const (
	statsPrefix = "STATS"
	initStats   = 0
)

func (s *Server) shortenURLController(request *ShortenURLRequest) (ShortenURLResponse, error) {
	hash := helpers.CreateUniqueHash()
	statKey := getStatsKey(hash)

	err := s.db.Set(hash, request.URL)
	if err != nil {
		return ShortenURLResponse{}, err
	}

	err = s.db.Set(statKey, strconv.Itoa(initStats))

	if err != nil {
		return ShortenURLResponse{}, err
	}

	returnURL := fmt.Sprintf("%s/%s", s.baseURL, hash)

	return ShortenURLResponse{ShortenedURL: returnURL}, nil
}

func (s *Server) getOriginalURLController(hash string) (url string, err error) {
	url, err = s.db.Get(hash)
	if err != nil {
		return "", err
	}

	statKey := getStatsKey(hash)
	_, err = s.db.Incr(statKey)

	if err != nil {
		return "", err
	}

	return url, err
}

func (s *Server) getStatsController(hash string) (res StatsResponse, err error) {
	url, err := s.db.Get(hash)
	if err != nil {
		return StatsResponse{}, err
	}

	statKey := getStatsKey(hash)

	hits, err := s.db.Get(statKey)
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
