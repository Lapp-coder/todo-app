package service

import (
	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/Lapp-coder/todo-app/internal/repository"
)

type TodoListService struct {
	repos repository.TodoList
}

func NewTodoListService(repos repository.TodoList) *TodoListService {
	return &TodoListService{repos: repos}
}

func (s TodoListService) Create(userID int, list model.TodoList) (int, error) {
	listID, err := s.repos.Create(userID, list)
	if err != nil {
		return 0, ErrFailedToCreateList
	}

	cacheList[userID] = listID

	return listID, nil
}

func (s TodoListService) GetAll(userID int) ([]model.TodoList, error) {
	lists, err := s.repos.GetAll(userID)
	if err != nil {
		return nil, ErrFailedToGetAllLists
	}

	return lists, nil
}

func (s TodoListService) GetByID(userID, listID int) (model.TodoList, error) {
	list, err := s.repos.GetByID(userID, listID)
	if err != nil {
		return model.TodoList{}, ErrFailedToGetListByID
	}

	return list, err
}

func (s TodoListService) Update(userID, listID int, update model.UpdateTodoList) error {
	if cacheList[userID] != listID {
		if _, err := s.GetByID(userID, listID); err != nil {
			return err
		}

		cacheList[userID] = listID
	}

	if err := s.repos.Update(userID, update); err != nil {
		return ErrFailedToUpdateList
	}

	return nil
}

func (s TodoListService) Delete(userID, listID int) error {
	if cacheList[userID] != listID {
		if _, err := s.repos.GetByID(userID, listID); err != nil {
			return err
		}
	}

	delete(cacheList, userID)
	delete(cacheItem, userID)

	if err := s.repos.Delete(listID); err != nil {
		return ErrFailedToDeleteList
	}

	return nil
}
