package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lapp-coder/todo-app/internal/app/handler"
	"github.com/Lapp-coder/todo-app/internal/app/repository"
	"github.com/Lapp-coder/todo-app/internal/app/server"
	"github.com/Lapp-coder/todo-app/internal/app/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Todo app API
// @version 1.0
// @description API server for todo list application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Завершение приложения при ошибке в инициализации конфигурационного файла
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializate a config file: %s", err.Error())
	}

	// Завершение приложения при ошибке в загрузке переменных окружения
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	// Инициализация бд postgres и завершение приложения при ошибке в подключении к бд
	db, err := repository.NewPostgresDB(repository.ConfigConnect{
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

	// Инициализация зависимостей
	repositories := repository.NewRepository(db)
	services := service.NewService(repositories)
	handlers := handler.NewHandler(services)

	// Инициализация http сервера
	srv := server.NewServer(server.Config{
		Host:           viper.GetString("server.host"),
		Port:           viper.GetString("server.port"),
		MaxHeaderBytes: viper.GetInt("server.maxHeaderBytes"),
		Handler:        handlers.InitRoutes(),
	})

	// Запуск сервера в go-рутине (для плавной остановки сервера)
	go func() {
		if err = srv.Run(); err != nil {
			logrus.Errorf("error occurred while starting the server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp started")

	// Ожидание на получение любого из сигналов из двух сигналов (SIGTERM, SIGINT) от системы для продолжения выполнения функции main()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp shutting down")

	// Плавная остановка сервера
	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}

	// Закрытие соединения с базой данных
	if err = db.Close(); err != nil {
		logrus.Errorf("error occurred when closing the connection to the database: %s", err.Error())
	}
}

// Инициализация конфигурационного файла
func initConfig() error {
	viper.AddConfigPath("configs/")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
