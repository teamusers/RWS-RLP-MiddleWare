package utils

import "github.com/gin-gonic/gin"

func RespondJSON(c *gin.Context, status int, message string, data any) {
	response := gin.H{
		"message": message,
	}

	if data != nil {
		response["data"] = data
	}

	c.JSON(status, response)
}
