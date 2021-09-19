package server

import "19shubham11/url-shortener/cmd/pkg/store"

type Server struct {
	db      store.Store
	baseURL string
}

func NewServer(db store.Store, baseURL string) *Server {
	return &Server{db: db, baseURL: baseURL}
}
