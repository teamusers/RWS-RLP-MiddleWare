package home

import (
	"context"
	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/codes"
	model "lbe/models"
	"lbe/system"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUsers handles GET /users
// If a user with the provided email already exists, it returns an error that the email not exists.
// If no user is found, it continues to generate an OTP.
func Login(c *gin.Context) {
	var req requests.LoginRequest
	// Bind the JSON payload to LoginRequest struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid email is required in the request body"})
		return
	}
	email := req.Email

	// Get a database handle.
	db := system.GetDb()

	// Attempt to find a user by email.
	var user model.User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		// If no user is found, return an error.
		if err == gorm.ErrRecordNotFound {
			resp := responses.APIResponse{
				Message: "email not found",
				Data: responses.LoginResponse{
					OTP:               "",
					ExpireIn:          0,
					LoginSessionToken: "",
					LoginExpireIn:     0,
				},
			}
			c.JSON(codes.CODE_EMAIL_NOTFOUND, resp)
			return

		}
		// For any other errors, return an internal server error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If the user exists, generate OTP using the service.
	otpService := services.NewOTPService()
	ctx := context.Background()
	otpResp, err := otpService.GenerateOTP(ctx, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	//Call send email services

	//Call IDP to return login session token & login expirein
	loginSessionToken := "idpSessionToken-ToBeProvide"
	login_expireIn := 30 * time.Minute //to be provide later
	loginExpireInSeconds := int64(login_expireIn.Seconds())

	resp := responses.APIResponse{
		Message: "email found",
		Data: responses.LoginResponse{
			OTP:               otpResp.OTP,
			ExpireIn:          otpResp.ExpiresAt,
			LoginSessionToken: loginSessionToken,
			LoginExpireIn:     loginExpireInSeconds,
		},
	}
	c.JSON(http.StatusOK, resp)
}
