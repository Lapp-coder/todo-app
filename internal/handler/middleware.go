package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	userCtx     = "userID"
	indexBearer = 0
	indexToken  = 1
)

func (h Handler) userAuthentication(ctx *gin.Context) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		respondError(ctx, http.StatusUnauthorized, errEmptyAuthHeader)
		return
	}

	headerParts := strings.Split(header, " ")
	bearerPart := headerParts[indexBearer]
	token := headerParts[indexToken]
	if len(headerParts) != 2 || bearerPart != "Bearer" {
		respondError(ctx, http.StatusUnauthorized, errInvalidAuthHeader)
		return
	}

	if token == "" {
		respondError(ctx, http.StatusUnauthorized, errEmptyToken)
		return
	}

	userID, err := h.service.Authorization.ParseToken(token)
	if err != nil {
		respondError(ctx, http.StatusUnauthorized, err)
		return
	}

	ctx.Set(userCtx, userID)
}

func (h Handler) getUserID(ctx *gin.Context) int {
	v, ok := ctx.Get(userCtx)
	if !ok {
		return 0
	}

	id, ok := v.(int)
	if !ok {
		return 0
	}

	return id
}
