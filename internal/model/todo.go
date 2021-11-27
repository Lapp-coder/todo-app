package model

type TodoList struct {
	ID             int    `json:"id" db:"id"`
	UserID         int    `json:"user_id" db:"user_id"`
	Title          string `json:"title" db:"title"`
	Description    string `json:"description" db:"description"`
	CompletionDate string `json:"completion_date" db:"completion_date"`
}

type TodoItem struct {
	ID             int    `json:"id" db:"id"`
	ListID         int    `json:"list_id" db:"list_id"`
	Title          string `json:"title" db:"title"`
	Description    string `json:"description" db:"description"`
	CompletionDate string `json:"completion_date" db:"completion_date"`
	Done           bool   `json:"done" db:"done"`
}
