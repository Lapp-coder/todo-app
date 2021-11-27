package postgres

import (
	"fmt"

	"github.com/Lapp-coder/todo-app/internal/config"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable     string = "users"
	todoListsTable string = "todo_lists"
	todoItemsTable string = "todo_items"
)

func NewDB(cfg config.PostgresDB) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	return db, nil
}
