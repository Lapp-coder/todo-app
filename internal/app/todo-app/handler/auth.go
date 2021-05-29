package handler

import (
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// signUp godoc
// @Summary Sign up
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body swagger.SignUpRequest true "Account info"
// @Success 201 {integer} integer "User id"
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /auth/sign-up [post]
func (h Handler) signUp(c *gin.Context) {
	req := struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if err := c.BindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.service.Authorization.CreateUser(model.User{Name: req.Name, Email: req.Email, Password: req.Password})
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	respond(c, http.StatusCreated, map[string]interface{}{
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
// @Param input body swagger.SignInRequest true "Credentials"
// @Success 201 {string} string "Token"
// @Failure 400 {object} swagger.ErrorResponse
// @Failure 404 {object} swagger.ErrorResponse
// @Failure 500 {object} swagger.ErrorResponse
// @Failure default {object} swagger.ErrorResponse
// @Router /auth/sign-in [post]
func (h Handler) signIn(c *gin.Context) {
	req := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if err := c.BindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	token, err := h.service.Authorization.GenerateToken(req.Email, req.Password)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	respond(c, http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
