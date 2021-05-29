package model

import validation "github.com/go-ozzo/ozzo-validation"

type TodoList struct {
	Id          int    `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
}

type TodoItem struct {
	Id          int    `db:"id"`
	ListId      int    `db:"list_id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Done        bool   `db:"done"`
}

func (l *TodoList) Validate() error {
	var fields []*validation.FieldRules

	if &l.Title != nil {
		fields = append(fields, validation.Field(&l.Title, validation.Length(2, 40)))
	}

	if &l.Description != nil {
		fields = append(fields, validation.Field(&l.Description, validation.Length(2, 100)))
	}

	return validation.ValidateStruct(
		l,
		fields...,
	)
}

func (i *TodoItem) Validate() error {
	var fields []*validation.FieldRules

	if &i.Title != nil {
		fields = append(fields, validation.Field(&i.Title, validation.Length(2, 40)))
	}

	if &i.Description != nil {
		fields = append(fields, validation.Field(&i.Description, validation.Length(2, 100)))
	}

	return validation.ValidateStruct(
		i,
		fields...,
	)
}
