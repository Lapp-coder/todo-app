package model

// User
type SignUp struct {
	Name     string `json:"name" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email,min=1,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

type SignIn struct {
	Email    string `json:"email" binding:"required,email,min=1,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

// List
type CreateTodoList struct {
	Title          string `json:"title" binding:"required,min=3,max=30"`
	Description    string `json:"description" binding:"max=50"`
	CompletionDate string `json:"completion_date"`
}

type UpdateTodoList struct {
	Title          *string `json:"title"`
	Description    *string `json:"description"`
	CompletionDate *string `json:"completion_date"`
}

func (l UpdateTodoList) IsNilAllFields() bool {
	return l.Title == nil && l.Description == nil && l.CompletionDate == nil
}

// Item
type CreateTodoItem struct {
	Title          string `json:"title" binding:"required,min=3,max=30"`
	Description    string `json:"description" binding:"max=50"`
	CompletionDate string `json:"completion_date"`
	Done           bool   `json:"done"`
}

type UpdateTodoItem struct {
	Title          *string `json:"title"`
	Description    *string `json:"description"`
	CompletionDate *string `json:"completion_date"`
	Done           *bool   `json:"done"`
}

func (i UpdateTodoItem) IsNilAllFields() bool {
	return i.Title == nil && i.Description == nil && i.CompletionDate == nil && i.Done == nil
}
