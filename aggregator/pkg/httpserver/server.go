package httpserver

import (
	"context"
	"net/http"
)

type Server struct {
	server *http.Server
}

func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler: handler,
	}

	server := &Server{
		server: httpServer,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
