package rest

import (
	"context"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewHTTPServer(address string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: handler,
		},
	}
}

func (s *Server) UP() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
