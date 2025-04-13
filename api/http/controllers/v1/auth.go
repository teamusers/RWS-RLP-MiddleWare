package home

import (
	"crypto/hmac"
	"fmt"
	"lbe/codes"
	model "lbe/models"
	"lbe/system"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// Adjust the import path based on your project structure and module name.
	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/api/http/services"

	"lbe/api/interceptor"
)

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
		resp := responses.ErrorResponse{
			Error: "Method Not Allowed",
		}
		c.JSON(http.StatusMethodNotAllowed, resp)
		return
	}

	// Check the Content-Type header.
	if c.GetHeader("Content-Type") != "application/json" {
		resp := responses.ErrorResponse{
			Error: "Content-Type must be application/json",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Retrieve the AppID from header.
	appID := c.GetHeader("AppID")
	if appID == "" {
		resp := responses.APIResponse{
			Message: "invalid appid",
			Data: responses.AuthResponse{
				AccessToken: "",
			},
		}
		c.JSON(codes.CODE_INVALID_APPID, resp)
		return
	}

	// Decode the JSON body.
	var req requests.AuthRequest
	if err := c.BindJSON(&req); err != nil {
		resp := responses.ErrorResponse{
			Error: "Invalid JSON body",
		}
		c.JSON(http.StatusMethodNotAllowed, resp)
		return
	}

	db := system.GetDb()
	// Look up the secret key associated with the AppID.
	secretKey, err := getSecretKey(db, appID)

	if err != nil || secretKey == "" {
		resp := responses.APIResponse{
			Message: "invalid appid",
			Data: responses.AuthResponse{
				AccessToken: "",
			},
		}
		c.JSON(codes.CODE_INVALID_APPID, resp)
		return
	}

	authReq, err := services.GenerateSignatureWithParams(appID, req.Nonce, req.Timestamp, secretKey)

	if err != nil {
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		// Handle the error, for example, send a JSON error response.
		return
	}
	expectedSignature := authReq.Signature

	// Compare the computed signature with the provided signature.
	if !hmac.Equal([]byte(expectedSignature), []byte(req.Signature)) {
		resp := responses.APIResponse{
			Message: "invalid signature",
			Data: responses.AuthResponse{
				AccessToken: "",
			},
		}
		c.JSON(codes.CODE_INVALID_SIGNATURE, resp)
		return
	}

	// Call the exported GenerateToken function from the middleware package.
	token, err := interceptor.GenerateToken(appID)
	if err != nil {
		resp := responses.ErrorResponse{
			Error: "Failed to generate token",
		}
		c.JSON(http.StatusMethodNotAllowed, resp)
		return
	}

	resp := responses.APIResponse{
		Message: "token successfully generated",
		Data: responses.AuthResponse{
			AccessToken: token,
		},
	}
	c.JSON(http.StatusOK, resp)
}
