package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func respond(c *gin.Context, statusCode int, v interface{}) {
	if err, ok := v.(error); ok {
		logrus.Errorf("%s, %s: %d, %s", c.ClientIP(), c.Request.UserAgent(), statusCode, err.Error())

		c.AbortWithStatusJSON(statusCode, map[string]string{
			"error": err.Error(),
		})

		return
	}

	c.JSON(statusCode, v)
}

func respondError(c *gin.Context, statusCode int, err error) {
	respond(c, statusCode, err)
}
