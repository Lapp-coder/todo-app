package postgres

import (
	"fmt"
	"strings"

	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/jmoiron/sqlx"
)

type TodoListRepository struct {
	db *sqlx.DB
}

func NewTodoListRepository(db *sqlx.DB) *TodoListRepository {
	return &TodoListRepository{db: db}
}

func (r *TodoListRepository) Create(userID int, list model.TodoList) (int, error) {
	var fields = make([]string, 0)
	var values = make([]interface{}, 0)
	var placeHolderID int
	var placeHolderIDs = make([]string, 0)

	fields = append(fields, "user_id")
	values = append(values, userID)
	placeHolderID++
	placeHolderIDs = append(placeHolderIDs, fmt.Sprintf("$%d", placeHolderID))

	if list.Title != "" {
		fields = append(fields, "title")
		values = append(values, list.Title)
		placeHolderID++
		placeHolderIDs = append(placeHolderIDs, fmt.Sprintf("$%d", placeHolderID))
	}

	if list.Description != "" {
		fields = append(fields, "description")
		values = append(values, list.Description)
		placeHolderID++
		placeHolderIDs = append(placeHolderIDs, fmt.Sprintf("$%d", placeHolderID))
	}

	if list.CompletionDate != "" {
		fields = append(fields, "completion_date")
		values = append(values, list.CompletionDate)
		placeHolderID++
		placeHolderIDs = append(placeHolderIDs, fmt.Sprintf("$%d", placeHolderID))
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		todoListsTable, strings.Join(fields, ","), strings.Join(placeHolderIDs, ","),
	)
	if err := r.db.QueryRow(query, values...).Scan(&list.ID); err != nil {
		return 0, err
	}

	return list.ID, nil
}

func (r *TodoListRepository) GetAll(userID int) ([]model.TodoList, error) {
	var lists []model.TodoList

	query := fmt.Sprintf(
		"SELECT tl.id, tl.user_id, tl.title, tl.description, tl.completion_date FROM %s tl WHERE tl.user_id = $1", todoListsTable)
	if err := r.db.Select(&lists, query, userID); err != nil {
		return nil, err
	}

	return lists, nil
}

func (r *TodoListRepository) GetByID(userID, listID int) (model.TodoList, error) {
	var list model.TodoList

	query := fmt.Sprintf(
		"SELECT tl.id, tl.user_id, tl.title, tl.description, tl.completion_date FROM %s tl WHERE tl.user_id = $1 AND tl.list_id = $2", todoListsTable)
	if err := r.db.Get(&list, query, userID, listID); err != nil {
		return list, err
	}

	return list, nil
}

func (r *TodoListRepository) Update(listID int, update model.UpdateTodoList) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	placeHolderID := 1

	if update.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", placeHolderID))
		args = append(args, *update.Title)
		placeHolderID++
	}

	if update.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", placeHolderID))
		args = append(args, *update.Description)
		placeHolderID++
	}

	if update.CompletionDate != nil {
		setValues = append(setValues, fmt.Sprintf("completion_date=$%d", placeHolderID))
		args = append(args, *update.CompletionDate)
		placeHolderID++
	}

	args = append(args, listID)

	query := fmt.Sprintf("UPDATE %s tl SET %s WHERE tl.id = $%d", todoListsTable, strings.Join(setValues, ", "), placeHolderID)
	if _, err := r.db.Exec(query, args...); err != nil {
		return err
	}

	return nil
}

func (r *TodoListRepository) Delete(listID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query1 := fmt.Sprintf("DELETE FROM %s ti WHERE ti.list_id = $1", todoItemsTable)
	if _, err = tx.Exec(query1, listID); err != nil {
		tx.Rollback()
		return err
	}

	query2 := fmt.Sprintf("DELETE FROM %s tl WHERE tl.id = $1", todoListsTable)
	if _, err = tx.Exec(query2, listID); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
