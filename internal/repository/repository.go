package repository

import (
	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/Lapp-coder/todo-app/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

var _ Authorization = (*postgres.AuthRepository)(nil)
var _ TodoList = (*postgres.TodoListRepository)(nil)
var _ TodoItem = (*postgres.TodoItem)(nil)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(email string) (model.User, error)
}

type TodoList interface {
	Create(userID int, list model.TodoList) (int, error)
	GetAll(userID int) ([]model.TodoList, error)
	GetByID(userID, listID int) (model.TodoList, error)
	Update(listID int, update model.UpdateTodoList) error
	Delete(listID int) error
}

type TodoItem interface {
	Create(listID int, item model.TodoItem) (int, error)
	GetAll(listID int) ([]model.TodoItem, error)
	GetByID(userID, itemID int) (model.TodoItem, error)
	Update(itemID int, update model.UpdateTodoItem) error
	Delete(itemID int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthRepository(db),
		TodoList:      postgres.NewTodoListRepository(db),
		TodoItem:      postgres.NewTodoItemRepository(db),
	}
}
