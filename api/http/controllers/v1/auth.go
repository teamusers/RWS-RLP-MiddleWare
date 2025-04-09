package home

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	model "rlp-member-service/models"
	"rlp-member-service/system"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// Adjust the import path based on your project structure and module name.
	"rlp-member-service/api/interceptor"
)

// AuthRequest represents the expected JSON body for the authentication request.
type AuthRequest struct {
	Timestamp string `json:"Timestamp"`
	Nonce     string `json:"Nonce"`
	Signature string `json:"Signature"`
}

// getSecretKey is a dummy function to lookup the secret key using the AppID.
// In a real implementation, this might query a database or another secure store.
func getSecretKey(db *gorm.DB, appID string) (string, error) {
	var channel model.SysChannel
	if err := db.Where("app_id = ?", appID).First(&channel).Error; err != nil {
		return "", fmt.Errorf("failed to get secret key for appID %s: %w", appID, err)
	}
	return channel.AppKey, nil
}

// AuthHandler processes the GET /api/v1/auth endpoint.
func AuthHandler(c *gin.Context) {
	// (Optional) Check the request method.
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method Not Allowed"})
		return
	}

	// Check the Content-Type header.
	if c.GetHeader("Content-Type") != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be application/json"})
		return
	}

	// Retrieve the AppID from header.
	appID := c.GetHeader("AppID")
	if appID == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "invalid appid",
			"data": gin.H{
				"accessToken": nil,
			},
		})
		return
	}

	// Decode the JSON body.
	var req AuthRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	db := system.GetDb()
	// Look up the secret key associated with the AppID.
	secretKey, err := getSecretKey(db, appID)

	if err != nil || secretKey == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "invalid appid",
			"data": gin.H{
				"accessToken": nil,
			},
		})
		return
	}

	// Concatenate AppID, Timestamp, and Nonce to create the base string.
	baseString := appID + req.Timestamp + req.Nonce
	fmt.Println("Base String:", baseString)

	// Compute the HMAC-SHA256 signature.
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(baseString))
	expectedMAC := mac.Sum(nil)
	expectedSignature := hex.EncodeToString(expectedMAC)
	fmt.Println("Expected Signature:", expectedSignature)

	// Compare the computed signature with the provided signature.
	if !hmac.Equal([]byte(expectedSignature), []byte(req.Signature)) {
		c.JSON(http.StatusOK, gin.H{
			"message": "invalid signature",
			"data": gin.H{
				"accessToken": nil,
			},
		})
		return
	}

	// At this point the request is authenticated.
	// Call the exported GenerateToken function from the middleware package.
	token, err := interceptor.GenerateToken(appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the JWT token in the JSON response.
	c.JSON(http.StatusOK, gin.H{
		"message": "token successfully generated",
		"data": gin.H{
			"accessToken": token,
		},
	})
}
