package v1

import (
	"crypto/hmac"
	"fmt"
	"lbe/codes"
	"lbe/model"
	"lbe/system"
	"log"
	"net/http"

	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/api/http/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"lbe/api/interceptor"
)

func getSecretKey(db *gorm.DB, appID string) (string, error) {
	var channel model.SysChannel
	if err := db.Where("app_id = ?", appID).First(&channel).Error; err != nil {
		return "", fmt.Errorf("failed to get secret key for appID %s: %w", appID, err)
	}
	return channel.AppKey, nil
}

// AuthHandler godoc
// @Summary      Generate authentication token
// @Description  Validates AppID header and HMAC signature, then returns a JWT access token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        AppID     header    string               true   "Client system AppID"
// @Param        request   body      requests.AuthRequest true   "Authentication request payload"
// @Success      200       {object}  responses.APIResponse{data=responses.AuthResponse}
// @Failure      400       {object}  responses.APIResponse
// @Failure      401       {object}  responses.APIResponse
// @Failure      500       {object}  responses.APIResponse
// @Router       /api/v1/auth [post]
func AuthHandler(c *gin.Context) {

	// Retrieve the AppID from header.
	appID := c.GetHeader("AppID")
	if appID == "" {
		resp := responses.APIResponse{
			Message: "missing appid",
			Data:    responses.AuthResponse{},
		}
		c.JSON(codes.CODE_INVALID_APPID, resp)
		return
	}

	// Decode the JSON body.
	var req requests.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := responses.APIResponse{
			Message: "invalid json request body",
			Data:    responses.AuthResponse{},
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	db := system.GetDb()
	// Look up the secret key associated with the AppID.
	secretKey, err := getSecretKey(db, appID)

	if err != nil || secretKey == "" {
		resp := responses.APIResponse{
			Message: "invalid appid",
			Data:    responses.AuthResponse{},
		}
		c.JSON(codes.CODE_INVALID_APPID, resp)
		return
	}

	authReq, err := services.GenerateSignatureWithParams(appID, req.Nonce, req.Timestamp, secretKey)

	if err != nil {
		log.Printf("error encountered generating auth signature: %v", err)
		resp := responses.APIResponse{
			Message: "internal error",
			Data:    responses.AuthResponse{},
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// Compare the computed signature with the provided signature.
	if !hmac.Equal([]byte(authReq.Signature), []byte(req.Signature)) {
		resp := responses.APIResponse{
			Message: "invalid signature",
			Data:    responses.AuthResponse{},
		}
		c.JSON(codes.CODE_INVALID_SIGNATURE, resp)
		return
	}

	// Call the exported GenerateToken function from the middleware package.
	token, err := interceptor.GenerateToken(appID)
	if err != nil {
		log.Printf("error encountered generating token: %v", err)
		resp := responses.APIResponse{
			Message: "internal error",
			Data:    responses.AuthResponse{},
		}
		c.JSON(http.StatusInternalServerError, resp)
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
