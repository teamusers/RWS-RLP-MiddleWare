package interceptor

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// CustomClaims represents the custom JWT claims.
type CustomClaims struct {
	AppID string `json:"app_id"`
	jwt.StandardClaims
}

// jwtSecret holds the secret key used for signing tokens.
// It is initially set to a default value.
var jwtSecret = []byte("RLP-Version1")

// SetJWTSecret allows you to update the jwtSecret from configuration.
func SetJWTSecret(secret string) {
	if secret != "" {
		jwtSecret = []byte(secret)
	}
}

// GenerateToken creates a JWT for the provided AppID.
func GenerateToken(appID string) (string, error) {
	// Set token expiration time (e.g., 1 hour from now)
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := CustomClaims{
		AppID: appID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "rlp-member-service-api", // Replace with your app name or identifier.
		},
	}

	// Create a new token object specifying the signing method and claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string.
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}
