package httpserver

import (
	"context"
	"net"
	"net/http"
)

type Server struct {
	server *http.Server
	errCh  chan error
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

func (s *Server) Start(ctx context.Context) {
	s.server.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}

	go func() {
		s.errCh <- s.server.ListenAndServe()
		close(s.errCh)
	}()
}

func (s *Server) Err() <-chan error {
	return s.errCh
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
