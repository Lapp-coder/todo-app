package handler

import (
	"net/http"
	"strconv"

	_ "github.com/Lapp-coder/todo-app/docs/swagger"
	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/gin-gonic/gin"
)

// createList godoc
// @Summary Create list
// @Security ApiKeyAuth
// @Tags lists
// @Description create list
// @ID create-list
// @Accept json
// @Produce json
// @Param input body model.CreateTodoList true "List info"
// @Success 201 {integer} integer "List id"
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/lists/ [post]
func (h Handler) createList(ctx *gin.Context) {
	userID := h.getUserID(ctx)
	if userID == 0 {
		respondError(ctx, http.StatusInternalServerError, errFailedToGetUserID)
		return
	}

	var req model.CreateTodoList
	if err := ctx.BindJSON(&req); err != nil {
		respondError(ctx, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	list := model.TodoList{Title: req.Title, Description: req.Description, CompletionDate: req.CompletionDate}
	listID, err := h.service.TodoList.Create(userID, list)
	if err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	respond(ctx, http.StatusCreated, gin.H{
		"list id": listID,
	})
}

// getAllLists godoc
// @Summary Get all lists
// @Security ApiKeyAuth
// @Tags lists
// @Description get all lists
// @ID get-all-lists
// @Produce json
// @Success 200 {object} swagger.GetAllListsResponse
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/lists/ [get]
func (h Handler) getAllLists(ctx *gin.Context) {
	userID := h.getUserID(ctx)
	if userID == 0 {
		respondError(ctx, http.StatusInternalServerError, errFailedToGetUserID)
		return
	}

	lists, err := h.service.TodoList.GetAll(userID)
	if err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	respond(ctx, http.StatusOK, gin.H{
		"lists": lists,
	})
}

// getListByID godoc
// @Summary Get list by id
// @Security ApiKeyAuth
// @Tags lists
// @Description get list by id
// @ID get-list-by-id
// @Produce json
// @Param id path int true "List id"
// @Success 200 {object} swagger.GetListByIDResponse
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/lists/{id} [get]
func (h Handler) getListByID(ctx *gin.Context) {
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

	list, err := h.service.TodoList.GetByID(userID, listID)
	if err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	respond(ctx, http.StatusOK, gin.H{
		"list": list,
	})
}

// updateList godoc
// @Summary Update list
// @Security ApiKeyAuth
// @Tags lists
// @Description update list by id
// @ID update-list
// @Accept json
// @Produce json
// @Param id path int true "List id"
// @Param input body model.UpdateTodoList true "Update values"
// @Success 200 {string} string "Result"
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/lists/{id} [put]
func (h Handler) updateList(ctx *gin.Context) {
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

	var req model.UpdateTodoList
	if err = ctx.BindJSON(&req); err != nil {
		respondError(ctx, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	if req.IsNilAllFields() {
		respondError(ctx, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	if err = h.service.TodoList.Update(userID, listID, req); err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	respond(ctx, http.StatusOK, gin.H{
		"result": "the list update was successful",
	})
}

// deleteList godoc
// @Summary Delete list
// @Security ApiKeyAuth
// @Tags lists
// @Description delete list by id
// @ID delete-list
// @Produce json
// @Param id path int true "List id"
// @Success 200 {string} string "Result"
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/lists/{id} [delete]
func (h Handler) deleteList(ctx *gin.Context) {
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

	if err = h.service.TodoList.Delete(userID, listID); err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	respond(ctx, http.StatusOK, gin.H{
		"result": "the list deletion was successful",
	})
}
