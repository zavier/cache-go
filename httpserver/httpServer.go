package httpserver

import (
	"../cache"
	"net/http"
)

type Server struct {
	cache.Cache
}

func (s *Server) Listen() {
	http.Handle("/status", s.statusHandler())
	http.ListenAndServe(":12345", nil)
}

func (s *Server) statusHandler() http.Handler {
	return &statusHandler{s}
}

func New(c cache.Cache) *Server {
	return &Server{c}
}
