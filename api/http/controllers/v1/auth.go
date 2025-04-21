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
// @Param        AppID     header    string               true   "Client system AppID" default(app1234)
// @Param        request   body      requests.AuthRequest true   "Authentication request payload"
// @Success      200       {object}  responses.AuthSuccessResponse "JWT access token returned successfully"
// @Failure      400       {object}  responses.ErrorResponse   "Malformed JSON in request body"
// @Failure      401       {object}  responses.ErrorResponse          "AppID header is missing"
// @Failure      401       {object}  responses.ErrorResponse          "AppID not recognized or unauthorized"
// @Failure      401       {object}  responses.ErrorResponse      "HMAC signature mismatch"
// @Failure      500       {object}  responses.ErrorResponse             "Unexpected server error"
// @Router       /auth [post]
func AuthHandler(c *gin.Context) {

	// Retrieve the AppID from header.
	appID := c.GetHeader("AppID")
	if appID == "" {
		c.JSON(http.StatusUnauthorized, responses.MissingAppIdErrorResponse())
		return
	}

	// Decode the JSON body.
	var req requests.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	db := system.GetDb()
	// Look up the secret key associated with the AppID.
	secretKey, err := getSecretKey(db, appID)

	if err != nil || secretKey == "" {
		c.JSON(http.StatusUnauthorized, responses.InvalidAppIdErrorResponse())
		return
	}

	authReq, err := services.GenerateSignatureWithParams(appID, req.Nonce, req.Timestamp, secretKey)

	if err != nil {
		log.Printf("error encountered generating auth signature: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	// Compare the computed signature with the provided signature.
	if !hmac.Equal([]byte(authReq.Signature), []byte(req.Signature)) {
		c.JSON(http.StatusUnauthorized, responses.InvalidSignatureErrorResponse())
		return
	}

	// Call the exported GenerateToken function from the middleware package.
	token, err := interceptor.GenerateToken(appID)
	if err != nil {
		log.Printf("error encountered generating token: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	resp := responses.ApiResponse[responses.AuthResponseData]{
		Code:    codes.SUCCESSFUL,
		Message: "token successfully generated",
		Data: responses.AuthResponseData{
			AccessToken: token,
		},
	}
	c.JSON(http.StatusOK, resp)
}
