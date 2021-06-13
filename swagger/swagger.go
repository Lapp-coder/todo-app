package swagger

import "github.com/Lapp-coder/todo-app/internal/app/todo-app/model"

type ErrorResponse struct {
	Error string `json:"error"`
}

type GetAllListsResponse struct {
	Lists []model.TodoList `json:"lists"`
}

type GetListByIdResponse struct {
	List model.TodoList `json:"list"`
}

type GetAllItemsResponse struct {
	Items []model.TodoItem `json:"items"`
}

type GetItemByIdResponse struct {
	Item model.TodoItem `json:"item"`
}
