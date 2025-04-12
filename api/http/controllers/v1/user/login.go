package home

import (
	"context"
	"lbe/api/http/requests"
	"lbe/api/http/services"
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
			c.JSON(201, gin.H{
				"message": "email not found",
				"data": gin.H{
					"otp":               nil,
					"otp_expireIn":      nil,
					"loginSessionToken": nil,
					"login_expireIn":    nil,
				},
			})
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

	// Return the response with the custom JSON format.
	c.JSON(http.StatusOK, gin.H{
		"message": "email found",
		"data": gin.H{
			"otp":               otpResp.OTP,
			"otp_expireIn":      otpResp.ExpiresAt,
			"loginSessionToken": loginSessionToken,
			"login_expireIn":    login_expireIn,
		},
	})

}
