package server

import (
	"fmt"
	"strconv"

	"19shubham11/url-shortener/cmd/internal/helpers"
	"19shubham11/url-shortener/cmd/pkg/redis"
)

const (
	statsPrefix = "STATS"
	initStats   = 0
)

func (s *Server) shortenURLController(request *ShortenURLRequest) (ShortenURLResponse, error) {
	hash := helpers.CreateUniqueHash()
	statKey := getStatsKey(hash)

	set := map[string]string{
		hash:    request.URL,
		statKey: strconv.Itoa(initStats),
	}

	err := s.db.Mset(set)

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
	statKey := getStatsKey(hash)
	keys := []string{hash, statKey}

	values, err := s.db.Mget(keys)

	if err != nil {
		return StatsResponse{}, err
	}

	url := values[0]
	hits := values[1]

	if hits == "" {
		return StatsResponse{}, redis.ErrorNotFound
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
