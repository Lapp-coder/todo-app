package request

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

type CreateTodoItem struct {
	Title          string `json:"title" binding:"required"`
	Description    string `json:"description"`
	CompletionDate string `json:"completion_date"`
	Done           bool   `json:"done"`
}

type UpdateTodoItem struct {
	Title          *string `json:"title"`
	Description    *string `json:"description"`
	CompletionDate *string `json:"completion_date"`
	Done           *bool   `json:"done"`
}

func (ci *CreateTodoItem) Validate() error {
	var fields []*validation.FieldRules

	fields = append(fields, validation.Field(&ci.Title, validation.Length(2, 40)))

	if &ci.Description != nil {
		fields = append(fields, validation.Field(&ci.Description, validation.Length(2, 100)))
	}

	if ci.CompletionDate == "" {
		ci.CompletionDate = defaultTime
	}

	if err := parseCompletedDate(&ci.CompletionDate); err != nil {
		return err
	}

	return validation.ValidateStruct(
		ci,
		fields...,
	)
}

func (ui *UpdateTodoItem) Validate() error {
	var fields []*validation.FieldRules

	if ui.Title == nil && ui.Description == nil && ui.Done == nil && ui.CompletionDate == nil {
		return errors.New("update request has not values")
	}

	if ui.Title != nil {
		fields = append(fields, validation.Field(&ui.Title, validation.Length(2, 40)))
	}

	if ui.Description != nil {
		fields = append(fields, validation.Field(&ui.Description, validation.Length(2, 100)))
	}

	if ui.CompletionDate != nil {
		if err := parseCompletedDate(ui.CompletionDate); err != nil {
			return errors.New("invalid input body")
		}
	}

	if err := validation.ValidateStruct(ui, fields...); err != nil {
		return errors.New("invalid input body")
	}

	return nil
}
