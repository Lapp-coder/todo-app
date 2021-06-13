package service

import (
	"errors"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/model"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/repository"
)

type TodoItemService struct {
	repos     repository.TodoItem
	reposList repository.TodoList
}

func NewTodoItemService(repos *repository.Repository) *TodoItemService {
	return &TodoItemService{repos: repos.TodoItem, reposList: repos.TodoList}
}

func (s TodoItemService) Create(userId, listId int, item model.TodoItem) (int, error) {
	if cacheList[userId] != listId {
		if _, err := s.reposList.GetById(userId, listId); err != nil {
			return 0, errors.New("failed to create item")
		}

		cacheList[userId] = listId
	}

	itemId, err := s.repos.Create(listId, item)
	if err != nil {
		return 0, err
	}

	cacheItem[userId] = itemId

	return itemId, nil
}

func (s TodoItemService) GetAll(userId, listId int) ([]model.TodoItem, error) {
	if cacheList[userId] != listId {
		if _, err := s.reposList.GetById(userId, listId); err != nil {
			return nil, errors.New("failed to get all the items")
		}

		cacheList[userId] = listId
	}

	return s.repos.GetAll(listId)
}

func (s TodoItemService) GetById(userId, itemId int) (model.TodoItem, error) {
	return s.repos.GetById(userId, itemId)
}

func (s TodoItemService) Update(userId, itemId int, item model.TodoItem) error {
	if cacheItem[userId] != itemId {
		if _, err := s.repos.GetById(userId, itemId); err != nil {
			return errors.New("failed to update item")
		}

		cacheItem[userId] = itemId
	}

	return s.repos.Update(itemId, item)
}

func (s TodoItemService) Delete(userId, itemId int) error {
	if cacheItem[userId] != itemId {
		if _, err := s.repos.GetById(userId, itemId); err != nil {
			return errors.New("failed to update item")
		}
	}

	delete(cacheItem, userId)

	return s.repos.Delete(itemId)
}
