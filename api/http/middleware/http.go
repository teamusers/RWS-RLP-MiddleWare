package middleware

import (
	"lbe/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpClientMiddleware(client *http.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := utils.WithHttpClient(c.Request.Context(), client)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
