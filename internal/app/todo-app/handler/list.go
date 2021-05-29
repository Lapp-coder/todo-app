package handler

import (
	"errors"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// createList godoc
// @Summary Create list
// @Security ApiKeyAuth
// @Tags lists
// @Description create list
// @ID create-list
// @Accept json
// @Produce json
// @Param input body swagger.CreateListRequest true "List info"
// @Success 201 {integer} integer "List id"
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/lists/ [post]
func (h Handler) createList(c *gin.Context) {
	userId := h.getUserId(c)

	req := struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}{}

	if err := c.BindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	listId, err := h.service.TodoList.Create(userId, model.TodoList{Title: req.Title, Description: req.Description})
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	respond(c, http.StatusCreated, map[string]interface{}{
		"list id": listId,
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
func (h Handler) getAllLists(c *gin.Context) {
	userId := h.getUserId(c)

	lists, err := h.service.TodoList.GetAll(userId)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	respond(c, http.StatusOK, map[string]interface{}{
		"lists": lists,
	})
}

// getListById godoc
// @Summary Get list by id
// @Security ApiKeyAuth
// @Tags lists
// @Description get list by id
// @ID get-list-by-id
// @Produce json
// @Param id path int true "List id"
// @Success 200 {object} swagger.GetListByIdResponse
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/lists/{id} [get]
func (h Handler) getListById(c *gin.Context) {
	userId := h.getUserId(c)

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid the param"))
		return
	}

	list, err := h.service.TodoList.GetById(userId, listId)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	respond(c, http.StatusOK, map[string]interface{}{
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
// @Param input body swagger.UpdateListRequest true "Update values"
// @Success 200 {string} string "Result"
// @Failure 400,404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /api/lists/{id} [put]
func (h Handler) updateList(c *gin.Context) {
	userId := h.getUserId(c)

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid the param"))
		return
	}

	req := struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
	}{}

	if err = c.BindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	input := struct{
		Title *string
		Description *string
	}{
		req.Title,
		req.Description,
	}
	if err = h.service.TodoList.Update(userId, listId, input); err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	respond(c, http.StatusOK, map[string]string{
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
func (h Handler) deleteList(c *gin.Context) {
	userId := h.getUserId(c)

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid the param"))
		return
	}

	if err = h.service.TodoList.Delete(userId, listId); err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	respond(c, http.StatusOK, map[string]string{
		"result": "the list deletion was successful",
	})
}
