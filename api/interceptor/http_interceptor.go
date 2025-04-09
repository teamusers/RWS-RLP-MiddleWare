package interceptor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// HttpInterceptor is a Gin middleware that validates the JWT access token.
// It expects the token to be provided in the "Authorization" header in the format "Bearer <token>".
func HttpInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access token"})
			c.Abort()
			return
		}

		// Expect the header to be in the format "Bearer <token>".
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse the token using the same secret and claims used for generating the token.
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		// Validate the token and ensure the claims are correct.
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			// Optionally, store the claims in the context for use by subsequent handlers.
			c.Set("claims", claims)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}
