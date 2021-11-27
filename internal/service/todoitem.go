package service

import (
	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/Lapp-coder/todo-app/internal/repository"
)

type TodoItemService struct {
	repos     repository.TodoItem
	reposList repository.TodoList
}

func NewTodoItemService(repos *repository.Repository) *TodoItemService {
	return &TodoItemService{repos: repos.TodoItem, reposList: repos.TodoList}
}

func (s TodoItemService) Create(userID, listID int, item model.TodoItem) (int, error) {
	if cacheList[userID] != listID {
		if _, err := s.reposList.GetByID(userID, listID); err != nil {
			return 0, ErrFailedToCreateItem
		}

		cacheList[userID] = listID
	}

	itemID, err := s.repos.Create(listID, item)
	if err != nil {
		return 0, err
	}

	cacheItem[userID] = itemID

	return itemID, nil
}

func (s TodoItemService) GetAll(userID, listID int) ([]model.TodoItem, error) {
	if cacheList[userID] != listID {
		if _, err := s.reposList.GetByID(userID, listID); err != nil {
			return nil, err
		}

		cacheList[userID] = listID
	}

	items, err := s.repos.GetAll(listID)
	if err != nil {
		return nil, ErrFailedToGetAllItems
	}

	return items, nil
}

func (s TodoItemService) GetByID(userID, itemID int) (model.TodoItem, error) {
	item, err := s.repos.GetByID(userID, itemID)
	if err != nil {
		return model.TodoItem{}, ErrFailedToGetItemByID
	}

	return item, nil
}

func (s TodoItemService) Update(userID, itemID int, update model.UpdateTodoItem) error {
	if cacheItem[userID] != itemID {
		if _, err := s.repos.GetByID(userID, itemID); err != nil {
			return err
		}

		cacheItem[userID] = itemID
	}

	if err := s.repos.Update(itemID, update); err != nil {
		return ErrFailedToUpdateItem
	}

	return nil
}

func (s TodoItemService) Delete(userID, itemID int) error {
	if cacheItem[userID] != itemID {
		if _, err := s.repos.GetByID(userID, itemID); err != nil {
			return err
		}
	}

	delete(cacheItem, userID)

	if err := s.repos.Delete(itemID); err != nil {
		return ErrFailedToDeleteItem
	}

	return nil
}
