package main

type ShortenURLRequest struct {
	URL string `json:"url"`
}

type ShortenURLResponse struct {
	ShortenedURL string `json:"shortenedURL"`
}

type StatsResponse struct {
	URL  string `json:"url"`
	Hits int    `json:"hits"`
}
