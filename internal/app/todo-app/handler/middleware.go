package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h Handler) userAuthentication(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		respondError(c, http.StatusUnauthorized, errors.New("empty auth header"))
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		respondError(c, http.StatusUnauthorized, errors.New("invalid auth header"))
		return
	}

	if headerParts[0] != "Bearer" {
		respondError(c, http.StatusUnauthorized, errors.New("invalid auth header"))
		return
	}

	if headerParts[1] == "" {
		respondError(c, http.StatusUnauthorized, errors.New("token is empty"))
		return
	}

	userId, err := h.service.Authorization.ParseToken(headerParts[1])
	if err != nil {
		respondError(c, http.StatusInternalServerError, errors.New("failed to parse token"))
		return
	}

	c.Set("userId", userId)
}

func (h Handler) getUserId(c *gin.Context) int {
	v, _ := c.Get("userId")

	id, _ := v.(int)

	return id
}
