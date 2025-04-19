package http

import (
	"context"
	"net/http"
	"time"
)

const (
	readHeaderTimeout = 20 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = time.Minute
)

type Server struct {
	httpSrv http.Server
}

func New(address string, handler http.Handler) *Server {
	server := &Server{
		httpSrv: http.Server{
			Addr:              address,
			Handler:           handler,
			ReadTimeout:       readTimeout,
			ReadHeaderTimeout: readHeaderTimeout,
			WriteTimeout:      writeTimeout,
		},
	}

	return server
}

func (s *Server) ListenAndServe() error {
	return s.httpSrv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}
