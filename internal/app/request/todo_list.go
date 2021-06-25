package request

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
)

type CreateTodoList struct {
	Title          string `json:"title" binding:"required"`
	Description    string `json:"description"`
	CompletionDate string `json:"completion_date"`
}

type UpdateTodoList struct {
	Title          *string `json:"title"`
	Description    *string `json:"description"`
	CompletionDate *string `json:"completion_date"`
}

func (cl *CreateTodoList) Validate() error {
	var fields []*validation.FieldRules

	fields = append(fields, validation.Field(&cl.Title, validation.Length(2, 40)))

	if &cl.Description != nil {
		fields = append(fields, validation.Field(&cl.Description, validation.Length(2, 100)))
	}

	if cl.CompletionDate == "" {
		cl.CompletionDate = getTimeNow()
		return validation.ValidateStruct(
			cl,
			fields...,
		)
	}

	if err := parseCompletedDate(&cl.CompletionDate); err != nil {
		return err
	}

	return validation.ValidateStruct(
		cl,
		fields...,
	)
}

func (ul *UpdateTodoList) Validate() error {
	var fields []*validation.FieldRules

	if ul.Title == nil && ul.Description == nil {
		return errors.New("update request has not values")
	}

	if ul.Title != nil {
		fields = append(fields, validation.Field(&ul.Title, validation.Length(2, 40)))
	}

	if ul.Description != nil {
		fields = append(fields, validation.Field(&ul.Description, validation.Length(2, 100)))
	}

	if ul.CompletionDate == nil {
		timeNow := getTimeNow()
		ul.CompletionDate = &timeNow
		return validation.ValidateStruct(
			ul,
			fields...,
		)
	}

	if err := parseCompletedDate(ul.CompletionDate); err != nil {
		return err
	}

	if err := validation.ValidateStruct(ul, fields...); err != nil {
		return errors.New("invalid input body")
	}

	return nil
}
