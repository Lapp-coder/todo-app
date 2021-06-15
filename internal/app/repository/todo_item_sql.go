package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Lapp-coder/todo-app/internal/app/model"
	"github.com/jmoiron/sqlx"
)

type TodoItemSQL struct {
	db *sqlx.DB
}

func NewTodoItemSQL(db *sqlx.DB) *TodoItemSQL {
	return &TodoItemSQL{db: db}
}

func (r *TodoItemSQL) Create(listId int, item model.TodoItem) (int, error) {
	err := r.db.QueryRow(fmt.Sprintf(
		"INSERT INTO %s (list_id, title, description) VALUES ($1, $2, $3) RETURNING id", todoItemsTable),
		listId, item.Title, item.Description).Scan(&item.Id)
	if err != nil {
		return 0, errors.New("error occurred when creating item")
	}

	return item.Id, nil
}

func (r *TodoItemSQL) GetAll(listId int) ([]model.TodoItem, error) {
	var items []model.TodoItem

	err := r.db.Select(&items, fmt.Sprintf(
		`SELECT ti.id, ti.list_id, ti.title, ti.description, ti.done FROM %s ti 
				INNER JOIN %s tl ON tl.id = ti.list_id WHERE tl.id = $1`, todoItemsTable, todoListsTable), listId)
	if err != nil {
		return nil, errors.New("error occurred when getting all items")
	}

	return items, nil
}

func (r *TodoItemSQL) GetById(userId, itemId int) (model.TodoItem, error) {
	var item model.TodoItem

	err := r.db.Get(&item, fmt.Sprintf(
		`SELECT ti.id, ti.list_id, ti.title, ti.description, ti.done FROM %s ti
				INNER JOIN %s ul ON ul.list_id = ti.list_id WHERE ul.user_id = $1 AND ti.id = $2`, todoItemsTable, usersListTable),
		userId, itemId)
	if err != nil {
		return item, errors.New("failed to get the item")
	}

	return item, nil
}

func (r *TodoItemSQL) Update(itemId int, item model.TodoItem) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	placeHolderId := 1

	if &item.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", placeHolderId))
		args = append(args, item.Title)
		placeHolderId++
	}

	if &item.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", placeHolderId))
		args = append(args, item.Description)
		placeHolderId++
	}

	if &item.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", placeHolderId))
		args = append(args, item.Done)
		placeHolderId++
	}

	args = append(args, itemId)

	_, err := r.db.Exec(fmt.Sprintf(
		"UPDATE %s ti SET %s WHERE ti.id = $%d",
		todoItemsTable, strings.Join(setValues, ", "), placeHolderId),
		args...)
	if err != nil {
		return errors.New("error occurred when item the item")
	}

	return nil
}

func (r *TodoItemSQL) Delete(itemId int) error {
	if _, err := r.db.Exec(fmt.Sprintf("DELETE FROM %s ti WHERE ti.id = $1", todoItemsTable), itemId); err != nil {
		return errors.New("error occurred when delete the item")
	}

	return nil
}
