package handler

import (
	"errors"
	"net/http"

	"github.com/Lapp-coder/todo-app/internal/app/todo-app/model"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/request"
	"github.com/gin-gonic/gin"
)

// signUp godoc
// @Summary Sign up
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body request.SignUp true "Account info"
// @Success 201 {integer} integer "User id"
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /auth/sign-up [post]
func (h Handler) signUp(c *gin.Context) {
	var req request.SignUp
	if err := c.BindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	if err := req.Validate(); err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	id, err := h.service.Authorization.CreateUser(model.User{Name: req.Name, Email: req.Email, Password: req.Password})
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	respond(c, http.StatusCreated, gin.H{
		"id": id,
	})
}

// signIn godoc
// @Summary Sign in
// @Tags auth
// @Description login
// @ID login
// @Accept json
// @Produce json
// @Param input body request.SignIn true "Credentials"
// @Success 201 {string} string "Token"
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /auth/sign-in [post]
func (h Handler) signIn(c *gin.Context) {
	var req request.SignIn
	if err := c.BindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	if err := req.Validate(); err != nil {
		respondError(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	token, err := h.service.Authorization.GenerateToken(req.Email, req.Password)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	respond(c, http.StatusOK, gin.H{
		"token": token,
	})
}
