package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lapp-coder/todo-app/internal/config"
	"github.com/Lapp-coder/todo-app/internal/handler"
	"github.com/Lapp-coder/todo-app/internal/repository"
	"github.com/Lapp-coder/todo-app/internal/repository/postgres"
	"github.com/Lapp-coder/todo-app/internal/server"
	"github.com/Lapp-coder/todo-app/internal/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const configPath = "configs/"

// @title Todo app API
// @version 2.1
// @description API server for todo list application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.New(configPath)
	if err != nil {
		logrus.Fatalf("failed to initializate config file: %s", err.Error())
	}

	db, err := postgres.NewDB(cfg.PostgresDB)
	if err != nil {
		logrus.Fatalf("failed to initializate db: %s", err.Error())
	}

	repositories := repository.New(db)
	services := service.New(repositories, cfg.Service)
	handlers := handler.New(services)

	cfg.Handler = handlers.InitRoutes()
	srv := server.NewServer(cfg.Server)

	go func() {
		if err = srv.Run(); err != nil {
			logrus.Errorf("failed to starting server: %s", err.Error())
		}
	}()

	logrus.Info("todo-app started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Info("todo-app shutting down")

	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("failed to shutting down server: %s", err.Error())
	}

	if err = db.Close(); err != nil {
		logrus.Errorf("failed to close connection to database: %s", err.Error())
	}
}
