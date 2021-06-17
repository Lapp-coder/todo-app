package repository

import (
	"github.com/Lapp-coder/todo-app/internal/app/model"
	"github.com/Lapp-coder/todo-app/internal/app/request"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(email string) (model.User, error)
}

type TodoList interface {
	Create(userId int, list model.TodoList) (int, error)
	GetAll(userId int) ([]model.TodoList, error)
	GetById(userId, listId int) (model.TodoList, error)
	Update(listId int, update request.UpdateTodoList) error
	Delete(listId int) error
}

type TodoItem interface {
	Create(listId int, item model.TodoItem) (int, error)
	GetAll(listId int) ([]model.TodoItem, error)
	GetById(userId, itemId int) (model.TodoItem, error)
	Update(itemId int, update request.UpdateTodoItem) error
	Delete(itemId int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSQL(db),
		TodoList:      NewTodoListSQL(db),
		TodoItem:      NewTodoItemSQL(db),
	}
}
