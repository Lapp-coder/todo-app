package server

import (
	"context"
	"net/http"
	"time"
)

type Config struct {
	Host           string
	Port           string
	MaxHeaderBytes int
	Handler        http.Handler
}

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg Config) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           cfg.Host + ":" + cfg.Port,
			Handler:        cfg.Handler,
			MaxHeaderBytes: cfg.MaxHeaderBytes << 20, // 1 MB
			ReadTimeout:    time.Second * 10,
			WriteTimeout:   time.Second * 10,
		},
	}
}

func (s Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
