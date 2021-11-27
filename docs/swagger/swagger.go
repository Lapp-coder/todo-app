package swagger

import "github.com/Lapp-coder/todo-app/internal/model"

type ErrorResponse struct {
	Error string `json:"error"`
}

type GetAllListsResponse struct {
	Lists []model.TodoList `json:"lists"`
}

type GetListByIDResponse struct {
	List model.TodoList `json:"list"`
}

type GetAllItemsResponse struct {
	Items []model.TodoItem `json:"items"`
}

type GetItemByIDResponse struct {
	Item model.TodoItem `json:"item"`
}
