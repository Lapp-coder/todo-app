package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Lapp-coder/todo-app/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg config.Server) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           cfg.Host + ":" + cfg.Port,
			Handler:        cfg.Handler,
			MaxHeaderBytes: cfg.MaxHeaderBytes << 20, // MB
			ReadTimeout:    time.Second * time.Duration(cfg.ReadTimeout),
			WriteTimeout:   time.Second * time.Duration(cfg.WriteTimeout),
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
