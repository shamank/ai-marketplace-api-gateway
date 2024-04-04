package app

import (
	"context"
	"net/http"
	"time"
)

type HTTPServer struct {
	http *http.Server
}

func NewHTTPServer(addr string, handler http.Handler, timeout time.Duration) *HTTPServer {

	return &HTTPServer{
		http: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
		},
	}
}

func (s *HTTPServer) Run() error {
	return s.http.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
