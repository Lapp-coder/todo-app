package postgres

import (
	"fmt"
	"strings"

	"github.com/Lapp-coder/todo-app/internal/model"
	TodoItemRepository "github.com/jmoiron/sqlx"
)

type TodoItem struct {
	db *TodoItemRepository.DB
}

func NewTodoItemRepository(db *TodoItemRepository.DB) *TodoItem {
	return &TodoItem{db: db}
}

func (r *TodoItem) Create(listID int, item model.TodoItem) (int, error) {
	var fields = make([]string, 0)
	var values = make([]interface{}, 0)
	var placeHolderID int
	var placeHolderIDs = make([]string, 0)

	fields = append(fields, "list_id")
	values = append(values, listID)
	placeHolderID++
	placeHolderIDs = append(placeHolderIDs, fmt.Sprintf("$%d", placeHolderID))

	if item.Title != "" {
		fields = append(fields, "title")
		values = append(values, item.Title)
		placeHolderID++
		placeHolderIDs = append(placeHolderIDs, fmt.Sprintf("$%d", placeHolderID))
	}

	if item.Description != "" {
		fields = append(fields, "description")
		values = append(values, item.Description)
		placeHolderID++
		placeHolderIDs = append(placeHolderIDs, fmt.Sprintf("$%d", placeHolderID))
	}

	if item.CompletionDate != "" {
		fields = append(fields, "completion_date")
		values = append(values, item.CompletionDate)
		placeHolderID++
		placeHolderIDs = append(placeHolderIDs, fmt.Sprintf("$%d", placeHolderID))
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		todoItemsTable, strings.Join(fields, ","), strings.Join(placeHolderIDs, ","),
	)
	if err := r.db.QueryRow(query, values...).Scan(&item.ID); err != nil {
		return 0, err
	}

	return item.ID, nil
}

func (r *TodoItem) GetAll(listID int) ([]model.TodoItem, error) {
	var items []model.TodoItem

	query := fmt.Sprintf(
		`SELECT ti.id, ti.list_id, ti.title, ti.description, ti.completion_date, ti.done FROM %s ti 
			WHERE ti.list_id = $1`, todoItemsTable)
	if err := r.db.Select(&items, query, listID); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItem) GetByID(userID, itemID int) (model.TodoItem, error) {
	var item model.TodoItem

	query := fmt.Sprintf(
		`SELECT ti.id, ti.list_id, ti.title, ti.description, ti.completion_date, ti.done FROM %s ti
				INNER JOIN %s tl ON tl.id = ti.list_id WHERE tl.user_id = $1 AND ti.id = $2`, todoItemsTable, todoListsTable)
	if err := r.db.Get(&item, query, userID, itemID); err != nil {
		return model.TodoItem{}, err
	}

	return item, nil
}

func (r *TodoItem) Update(itemID int, update model.UpdateTodoItem) error {
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

	if update.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", placeHolderID))
		args = append(args, *update.Done)
		placeHolderID++
	}

	args = append(args, itemID)

	query := fmt.Sprintf(
		"UPDATE %s ti SET %s WHERE ti.id = $%d",
		todoItemsTable, strings.Join(setValues, ", "), placeHolderID)
	if _, err := r.db.Exec(query, args...); err != nil {
		return err
	}

	return nil
}

func (r *TodoItem) Delete(itemID int) error {
	query := fmt.Sprintf("DELETE FROM %s ti WHERE ti.id = $1", todoItemsTable)
	if _, err := r.db.Exec(query, itemID); err != nil {
		return err
	}

	return nil
}
