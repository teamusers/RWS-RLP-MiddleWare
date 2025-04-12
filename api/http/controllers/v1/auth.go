package home

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"rlp-middleware/config"
	model "rlp-middleware/models"
	"rlp-middleware/system"
	"rlp-middleware/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// Adjust the import path based on your project structure and module name.
	"rlp-middleware/api/http/requests"
	"rlp-middleware/api/http/responses"
	"rlp-middleware/api/http/services"
	"rlp-middleware/api/interceptor"
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
	var req requests.AuthRequest
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
	baseString := appID + *req.Timestamp + *req.Nonce
	fmt.Println("Base String:", baseString)

	// Compute the HMAC-SHA256 signature.
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(baseString))
	expectedMAC := mac.Sum(nil)
	expectedSignature := hex.EncodeToString(expectedMAC)
	fmt.Println("Expected Signature:", expectedSignature)

	// Compare the computed signature with the provided signature.
	if !hmac.Equal([]byte(expectedSignature), []byte(*req.Sign)) {
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

func InitiateLogin(c *gin.Context) {
	var req requests.InitiateLoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	// check email address existence
	// TODO: Call Member Service here

	if req.DebugMode != nil {
		log.Println("debug mode enabled for initate login")

		switch *req.DebugMode {
		case 2:
			log.Println("simulate email NOT found")
			utils.RespondJSON(c, http.StatusCreated, "email not found", responses.LoginResponse{})
			return
		case 3:
			log.Println("simulate network error")
			utils.RespondJSON(c, http.StatusBadRequest, "internal error", responses.LoginResponse{})
			return
		default:
			log.Println("simulate email found")
		}
	}

	// user found
	log.Println("valid email address found")

	// generate otp
	otpService := services.NewOTPService()
	otpResp, err := otpService.GenerateOTP(c, *req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate OTP"})
		return
	}

	// request login session token
	// TODO: Call Member Service here
	loginSessionToken := "sampleToken"
	loginSessionExpiry := "1298738172489"

	//send otp via email
	emailData := services.TemplateData{
		Email: *req.Email,
		OTP:   otpResp.OTP,
	}

	cfg := config.GetConfig()
	emailService := services.NewEmailService(&cfg.Smtp)
	if err := emailService.SendOtpEmail(*req.Email, emailData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to send email otp"})
		return
	}

	// Return response
	response := responses.LoginResponse{
		OTP:               &otpResp.OTP,
		ExpireIn:          &otpResp.ExpiresAt,
		LoginSessionToken: &loginSessionToken,
		LoginExpireIn:     &loginSessionExpiry,
	}

	utils.RespondJSON(c, http.StatusOK, "email found", response)
}
