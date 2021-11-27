package handler

import (
	"net/http"
	"strconv"

	_ "github.com/Lapp-coder/todo-app/docs/swagger"
	"github.com/Lapp-coder/todo-app/internal/model"
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
// @Param input body model.CreateTodoItem true "Item info"
// @Success 201 {integer} integer "Item id"
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/lists/{id}/items/ [post]
func (h Handler) createItem(ctx *gin.Context) {
	userID := h.getUserID(ctx)
	if userID == 0 {
		respondError(ctx, http.StatusInternalServerError, errFailedToGetUserID)
		return
	}

	listID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		respondError(ctx, http.StatusBadRequest, errInvalidParamID)
		return
	}

	var req model.CreateTodoItem
	if err = ctx.BindJSON(&req); err != nil {
		respondError(ctx, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	item := model.TodoItem{Title: req.Title, Description: req.Description, CompletionDate: req.CompletionDate, Done: req.Done}
	itemID, err := h.service.TodoItem.Create(userID, listID, item)
	if err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	respond(ctx, http.StatusCreated, gin.H{
		"item_id": itemID,
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
func (h Handler) getAllItems(ctx *gin.Context) {
	userID := h.getUserID(ctx)
	if userID == 0 {
		respondError(ctx, http.StatusInternalServerError, errFailedToGetUserID)
		return
	}

	listID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		respondError(ctx, http.StatusBadRequest, errInvalidParamID)
		return
	}

	items, err := h.service.TodoItem.GetAll(userID, listID)
	if err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	respond(ctx, http.StatusOK, gin.H{
		"items": items,
	})
}

// getItemByID godoc
// @Summary Get item by id
// @Security ApiKeyAuth
// @Tags items
// @Description get item by id
// @ID get-item-by-id
// @Produce json
// @Param id path int true "Item id"
// @Success 200 {object} swagger.GetItemByIDResponse
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/items/{id} [get]
func (h Handler) getItemByID(ctx *gin.Context) {
	userID := h.getUserID(ctx)
	if userID == 0 {
		respondError(ctx, http.StatusInternalServerError, errFailedToGetUserID)
		return
	}

	itemID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		respondError(ctx, http.StatusBadRequest, errInvalidParamID)
		return
	}

	item, err := h.service.TodoItem.GetByID(userID, itemID)
	if err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	respond(ctx, http.StatusOK, gin.H{
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
// @Param input body model.UpdateTodoItem true "Update values"
// @Success 200 {string} string "Result"
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/items/{id} [put]
func (h Handler) updateItem(ctx *gin.Context) {
	userID := h.getUserID(ctx)
	if userID == 0 {
		respondError(ctx, http.StatusInternalServerError, errFailedToGetUserID)
		return
	}

	itemID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		respondError(ctx, http.StatusBadRequest, errInvalidParamID)
		return
	}

	var req model.UpdateTodoItem
	if err = ctx.BindJSON(&req); err != nil {
		respondError(ctx, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	if req.IsNilAllFields() {
		respondError(ctx, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	if err = h.service.TodoItem.Update(userID, itemID, req); err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	respond(ctx, http.StatusOK, gin.H{
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
func (h Handler) deleteItem(ctx *gin.Context) {
	userID := h.getUserID(ctx)
	if userID == 0 {
		respondError(ctx, http.StatusInternalServerError, errFailedToGetUserID)
		return
	}

	itemID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		respondError(ctx, http.StatusBadRequest, errInvalidParamID)
		return
	}

	if err = h.service.TodoItem.Delete(userID, itemID); err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	respond(ctx, http.StatusOK, gin.H{
		"result": "the item deletion was successful",
	})
}
