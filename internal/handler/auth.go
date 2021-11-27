package handler

import (
	"net/http"

	_ "github.com/Lapp-coder/todo-app/docs/swagger"
	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/gin-gonic/gin"
)

// signUp godoc
// @Summary Sign up
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body model.SignUp true "Account info"
// @Success 201 {integer} integer "User id"
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /auth/sign-up [post]
func (h Handler) signUp(ctx *gin.Context) {
	var req model.SignUp
	if err := ctx.BindJSON(&req); err != nil {
		respondError(ctx, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	user := model.User{Name: req.Name, Email: req.Email, Password: req.Password}
	id, err := h.service.Authorization.CreateUser(user)
	if err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	respond(ctx, http.StatusCreated, gin.H{
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
// @Param input body model.SignIn true "Credentials"
// @Success 201 {string} string "Token"
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /auth/sign-in [post]
func (h Handler) signIn(ctx *gin.Context) {
	var req model.SignIn
	if err := ctx.BindJSON(&req); err != nil {
		respondError(ctx, http.StatusBadRequest, errInvalidInputBody)
		return
	}

	token, err := h.service.Authorization.GenerateToken(req.Email, req.Password)
	if err != nil {
		respondError(ctx, http.StatusInternalServerError, err)
		return
	}

	respond(ctx, http.StatusOK, gin.H{
		"token": token,
	})
}
