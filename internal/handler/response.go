package handler

import (
	"github.com/gin-gonic/gin"
)

func respond(ctx *gin.Context, statusCode int, value interface{}) {
	if err, ok := value.(error); ok {
		ctx.AbortWithStatusJSON(statusCode, map[string]string{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(statusCode, value)
}

func respondError(ctx *gin.Context, statusCode int, err error) {
	respond(ctx, statusCode, err)
}
