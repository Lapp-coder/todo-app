package model

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
