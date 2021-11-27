package service

import "errors"

var (
	ErrIncorrectEmailOrPassword = errors.New("incorrect email or password")
	ErrInvalidSigningMethod     = errors.New("invalid signing method")
	ErrFailedToCreateItem       = errors.New("failed to create item")
	ErrFailedToGetAllItems      = errors.New("failed to get all items")
	ErrFailedToGetItemByID      = errors.New("failed to get item by id")
	ErrFailedToUpdateItem       = errors.New("failed to update item")
	ErrFailedToDeleteItem       = errors.New("failed to delete item")
	ErrFailedToCreateList       = errors.New("failed to create list")
	ErrFailedToGetAllLists      = errors.New("failed to get all lists")
	ErrFailedToGetListByID      = errors.New("failed to get list by id")
	ErrFailedToUpdateList       = errors.New("failed to update list")
	ErrFailedToDeleteList       = errors.New("failed to delete list")
)
