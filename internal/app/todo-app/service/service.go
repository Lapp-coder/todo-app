package service

import (
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/model"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/repository"
)

var (
	cacheList = make(map[int]int)
	cacheItem = make(map[int]int)
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type TodoList interface {
	Create(userId int, list model.TodoList) (int, error)
	GetAll(userId int) ([]model.TodoList, error)
	GetById(userId, listId int) (model.TodoList, error)
	Update(userId, listId int, input struct {Title *string; Description *string}) error
	Delete(userId, listId int) error
}

type TodoItem interface {
	Create(userId, listId int, item model.TodoItem) (int, error)
	GetAll(userId, listId int) ([]model.TodoItem, error)
	GetById(userId, itemId int) (model.TodoItem, error)
	Update(userId, itemId int, input struct {Title *string; Description *string; Done *bool}) error
	Delete(userId, itemId int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos),
	}
}
