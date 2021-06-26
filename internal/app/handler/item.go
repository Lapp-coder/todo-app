package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Lapp-coder/todo-app/internal/app/model"
	"github.com/Lapp-coder/todo-app/internal/app/request"
	"github.com/gin-gonic/gin"
)

// createItem godoc
// @Summary Create item
// @Security ApiKeyAuth
// @Tags items
// @Description create item
// @ID create-item
// @Accept json
// @Produce json
// @Param id path int true "List id"
// @Param input body request.CreateTodoItem true "Item info"
// @Success 201 {integer} integer "Item id"
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/lists/{id}/items/ [post]
func (h Handler) createItem(c *gin.Context) {
	userId := h.getUserId(c)
	if userId == 0 {
		respondError(c, http.StatusInternalServerError, errors.New("failed to get user id"))
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid the param"))
		return
	}

	var req request.CreateTodoItem
	if err = c.BindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	if err = req.Validate(); err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	item := model.TodoItem{Title: req.Title, Description: req.Description, CompletionDate: req.CompletionDate, Done: req.Done}
	itemId, err := h.service.TodoItem.Create(userId, listId, item)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	respond(c, http.StatusCreated, gin.H{
		"item id": itemId,
	})
}

// getAllItems godoc
// @Summary Get all items
// @Security ApiKeyAuth
// @Tags items
// @Description get all items
// @ID get-all-items
// @Produce json
// @Param id path int true "List id"
// @Success 200 {object} swagger.GetAllItemsResponse
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/lists/{id}/items/ [get]
func (h Handler) getAllItems(c *gin.Context) {
	userId := h.getUserId(c)
	if userId == 0 {
		respondError(c, http.StatusInternalServerError, errors.New("failed to get user id"))
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid the param"))
		return
	}

	items, err := h.service.TodoItem.GetAll(userId, listId)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	respond(c, http.StatusOK, gin.H{
		"items": items,
	})
}

// getItemById godoc
// @Summary Get item by id
// @Security ApiKeyAuth
// @Tags items
// @Description get item by id
// @ID get-item-by-id
// @Produce json
// @Param id path int true "Item id"
// @Success 200 {object} swagger.GetItemByIdResponse
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/items/{id} [get]
func (h Handler) getItemById(c *gin.Context) {
	userId := h.getUserId(c)
	if userId == 0 {
		respondError(c, http.StatusInternalServerError, errors.New("failed to get user id"))
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid the param"))
		return
	}

	item, err := h.service.TodoItem.GetById(userId, itemId)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	respond(c, http.StatusOK, gin.H{
		"item": item,
	})
}

// updateItem godoc
// @Summary Update item
// @Security ApiKeyAuth
// @Tags items
// @Description update item by id
// @ID update-item
// @Accept json
// @Produce json
// @Param id path int true "Item id"
// @Param input body request.UpdateTodoItem true "Update values"
// @Success 200 {string} string "Result"
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/items/{id} [put]
func (h Handler) updateItem(c *gin.Context) {
	userId := h.getUserId(c)
	if userId == 0 {
		respondError(c, http.StatusInternalServerError, errors.New("failed to get user id"))
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid the param"))
		return
	}

	var req request.UpdateTodoItem
	if err = c.BindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	if err = req.Validate(); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	err = h.service.TodoItem.Update(userId, itemId, req)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	respond(c, http.StatusOK, gin.H{
		"result": "the item update was successful",
	})
}

// deleteItem godoc
// @Summary Delete item
// @Security ApiKeyAuth
// @Tags items
// @Description delete item by id
// @ID delete-item
// @Produce json
// @Param id path int true "Item id"
// @Success 200 {string} string "Result"
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/items/{id} [delete]
func (h Handler) deleteItem(c *gin.Context) {
	userId := h.getUserId(c)
	if userId == 0 {
		respondError(c, http.StatusInternalServerError, errors.New("failed to get user id"))
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid the param"))
		return
	}

	if err = h.service.TodoItem.Delete(userId, itemId); err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	respond(c, http.StatusOK, gin.H{
		"result": "the item deletion was successful",
	})
}
