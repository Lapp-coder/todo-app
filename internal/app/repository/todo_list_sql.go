package repository

import (
	"errors"
	"fmt"
	"github.com/Lapp-coder/todo-app/internal/app/model"
	"github.com/Lapp-coder/todo-app/internal/app/request"
	"github.com/jmoiron/sqlx"
	"strings"
)

type TodoListSQL struct {
	db *sqlx.DB
}

func NewTodoListSQL(db *sqlx.DB) *TodoListSQL {
	return &TodoListSQL{db: db}
}

func (r *TodoListSQL) Create(userId int, list model.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, errors.New("an error occurred when creating a list")
	}

	err = r.db.QueryRow(fmt.Sprintf(
		"INSERT INTO %s (title, description, completion_date) VALUES ($1, $2, $3) RETURNING id", todoListsTable),
		list.Title, list.Description, list.CompletionDate).Scan(&list.Id)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("an error occurred when creating a list")
	}

	_, err = r.db.Exec(fmt.Sprintf(
		"INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable), userId, list.Id)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("an error occurred when creating a list")
	}

	if err = tx.Commit(); err != nil {
		return 0, errors.New("an error occurred when creating a list")
	}

	return list.Id, nil
}

func (r *TodoListSQL) GetAll(userId int) ([]model.TodoList, error) {
	var lists []model.TodoList

	err := r.db.Select(&lists, fmt.Sprintf(
		"SELECT tl.id, tl.title, tl.description, tl.completion_date FROM %s tl INNER JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable), userId)
	if err != nil {
		return nil, errors.New("failed to get all the lists")
	}

	return lists, nil
}

func (r *TodoListSQL) GetById(userId, listId int) (model.TodoList, error) {
	var list model.TodoList

	err := r.db.Get(&list, fmt.Sprintf(
		"SELECT tl.id, tl.title, tl.description, tl.completion_date FROM %s tl INNER JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2",
		todoListsTable, usersListsTable), userId, listId)
	if err != nil {
		return list, errors.New("failed to get the list")
	}

	return list, nil
}

func (r *TodoListSQL) Update(listId int, update request.UpdateTodoList) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	placeHolderId := 1

	if update.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", placeHolderId))
		args = append(args, *update.Title)
		placeHolderId++
	}

	if update.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", placeHolderId))
		args = append(args, *update.Description)
		placeHolderId++
	}

	if update.CompletionDate != nil {
		setValues = append(setValues, fmt.Sprintf("completion_date=$%d", placeHolderId))
		args = append(args, *update.CompletionDate)
		placeHolderId++
	}

	args = append(args, listId)

	_, err := r.db.Exec(fmt.Sprintf(
		"UPDATE %s tl SET %s WHERE tl.id = $%d",
		todoListsTable, strings.Join(setValues, ", "), placeHolderId),
		args...)
	if err != nil {
		return errors.New("an error occurred when updating the list")
	}

	return nil
}

func (r *TodoListSQL) Delete(listId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return errors.New("an error occurred when deleting the list")
	}

	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s ti WHERE ti.list_id = $1", todoItemsTable), listId)
	if err != nil {
		tx.Rollback()
		return errors.New("an error occurred when deleting the list")
	}

	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s tl WHERE tl.id = $1", todoListsTable), listId)
	if err != nil {
		tx.Rollback()
		return errors.New("an error occurred when deleting the list")
	}

	_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s ul WHERE ul.list_id = $1", usersListsTable), listId)
	if err != nil {
		tx.Rollback()
		return errors.New("an error occurred when deleting the list")
	}

	if err = tx.Commit(); err != nil {
		return errors.New("an error occurred when deleting the list")
	}

	return nil
}
