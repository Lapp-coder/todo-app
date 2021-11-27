package service

import (
	"github.com/Lapp-coder/todo-app/internal/config"
	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/Lapp-coder/todo-app/internal/repository"
)

//go:generate mockgen -source=service.go --destination=mocks/mock.go

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
	Create(userID int, list model.TodoList) (int, error)
	GetAll(userID int) ([]model.TodoList, error)
	GetByID(userID, listID int) (model.TodoList, error)
	Update(userID, listID int, update model.UpdateTodoList) error
	Delete(userID, listID int) error
}

type TodoItem interface {
	Create(userID, listID int, item model.TodoItem) (int, error)
	GetAll(userID, listID int) ([]model.TodoItem, error)
	GetByID(userID, itemID int) (model.TodoItem, error)
	Update(userID, itemID int, update model.UpdateTodoItem) error
	Delete(userID, itemID int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func New(repos *repository.Repository, cfg config.Service) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, cfg),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos),
	}
}
