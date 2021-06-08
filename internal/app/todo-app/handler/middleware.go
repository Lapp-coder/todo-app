package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	userCtx string = "userId"
)

func (h Handler) userAuthentication(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		respondError(c, http.StatusUnauthorized, errors.New("empty auth header"))
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		respondError(c, http.StatusUnauthorized, errors.New("invalid auth header"))
		return
	}

	if headerParts[1] == "" {
		respondError(c, http.StatusUnauthorized, errors.New("token is empty"))
		return
	}

	userId, err := h.service.Authorization.ParseToken(headerParts[1])
	if err != nil {
		respondError(c, http.StatusUnauthorized, errors.New("failed to parse token"))
		return
	}

	c.Set(userCtx, userId)
}

func (h Handler) getUserId(c *gin.Context) int {
	v, ok := c.Get(userCtx)
	if !ok {
		return 0
	}

	id, ok := v.(int)
	if !ok {
		return 0
	}

	return id
}
