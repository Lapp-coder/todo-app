package main

import (
	"context"
	"flag"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/handler"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/repository"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/server"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "c", "./config", "path to config file")
}

func initConfig() error {
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

// @title Todo app API
// @version 1.0
// @description API server for todo list application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	flag.Parse()

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config file: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("db.db_name"),
		SSLMode:  viper.GetString("db.ssl_mode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initializate db: %s", err.Error())
	}

	repositories := repository.NewRepository(db)
	services := service.NewService(repositories)
	handlers := handler.NewHandler(services)

	srv := server.NewServer(server.Config{
		Host:           viper.GetString("server.host"),
		Port:           viper.GetString("server.port"),
		MaxHeaderBytes: viper.GetInt("server.maxHeaderBytes"),
		Handler:        handlers.InitRoutes(),
	})
	go func() {
		if err = srv.Run(); err != nil {
			logrus.Errorf("error occurred while starting the server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp shutting down")

	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}

	if err = db.Close(); err != nil {
		logrus.Errorf("error occurred when closing the connection to the database: %s", err.Error())
	}
}
