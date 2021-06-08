package swagger

import "github.com/Lapp-coder/todo-app/internal/app/todo-app/model"

type ErrorResponse struct {
	Error string `json:"error"`
}

type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateListRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateListRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type GetAllListsResponse struct {
	Lists []model.TodoList `json:"lists"`
}

type GetListByIdResponse struct {
	List model.TodoList `json:"list"`
}

type CreateItemRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateItemRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

type GetAllItemsResponse struct {
	Items []model.TodoItem `json:"items"`
}

type GetItemByIdResponse struct {
	Item model.TodoItem `json:"item"`
}
