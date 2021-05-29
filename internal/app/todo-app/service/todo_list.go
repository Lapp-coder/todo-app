package service

import (
	"errors"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/model"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/repository"
)

type TodoListService struct {
	repos repository.TodoList
}

func NewTodoListService(repos repository.TodoList) *TodoListService {
	return &TodoListService{repos: repos}
}

func (s TodoListService) Create(userId int, list model.TodoList) (int, error) {
	if err := list.Validate(); err != nil {
		return 0, err
	}

	listId, err := s.repos.Create(userId, list)
	if err != nil {
		return 0, err
	}

	cacheList[userId] = listId

	return listId, nil
}

func (s TodoListService) GetAll(userId int) ([]model.TodoList, error) {
	return s.repos.GetAll(userId)
}

func (s TodoListService) GetById(userId, listId int) (model.TodoList, error) {
	return s.repos.GetById(userId, listId)
}

func (s TodoListService) Update(userId, listId int, input struct {Title *string; Description *string}) error {
	if input.Title == nil && input.Description == nil {
		return errors.New("update request has not values")
	}

	var list model.TodoList
	if input.Title != nil {
		list.Title = *input.Title
	}
	if input.Description != nil {
		list.Description = *input.Description
	}

	if err := list.Validate(); err != nil {
		return err
	}

	if cacheList[userId] != listId {
		if _, err := s.GetById(userId, listId); err != nil {
			return errors.New("failed to update list")
		}

		cacheList[userId] = listId
	}

	return s.repos.Update(listId, input)
}

func (s TodoListService) Delete(userId, listId int) error {
	if cacheList[userId] != listId {
		if _, err := s.repos.GetById(userId, listId); err != nil {
			return errors.New("failed to delete list")
		}
	}

	delete(cacheList, userId)
	delete(cacheItem, userId)

	return s.repos.Delete(listId)
}
