package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable     string = "users"
	todoListsTable string = "todo_lists"
	todoItemsTable string = "todo_items"
	usersListTable string = "users_lists"
)

type ConfigConnect struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg ConfigConnect) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))

	return db, err
}
