package home

import (
	"context"
	"net/http"
	"time"

	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/codes"
	model "lbe/models"
	"lbe/system"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUser(c *gin.Context) {
	email := c.Param("email")
	signUpType := c.Param("sign_up_type")

	if email == "" || signUpType == "" {
		resp := responses.ErrorResponse{
			Error: "Valid email and sign_up_type are required as query parameters",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	switch signUpType {
	case "new":
		db := system.GetDb()
		// Attempt to find a user by email.
		var user model.User
		err := db.Preload("PhoneNumbers").Where("email = ?", email).First(&user).Error
		if err == nil {
			// User found: return an error indicating that the email already exists.
			resp := responses.APIResponse{
				Message: "email registered",
				Data: responses.SignUpResponse{
					OTP:      "",
					ExpireIn: 0,
				},
			}
			c.JSON(codes.CODE_EMAIL_REGISTERED, resp)
			return
		}
		if err != gorm.ErrRecordNotFound {
			// An unexpected error occurred during the query.
			resp := responses.ErrorResponse{
				Error: err.Error(),
			}
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		// Generate OTP using the service.
		otpService := services.NewOTPService()
		ctx := context.Background()
		otpResp, err := otpService.GenerateOTP(ctx, email)
		if err != nil {
			resp := responses.ErrorResponse{
				Error: "Failed to generate OTP",
			}
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		//send email
		resp := responses.APIResponse{
			Message: "email not registered",
			Data: responses.SignUpResponse{
				OTP:      otpResp.OTP,
				ExpireIn: otpResp.ExpiresAt,
			},
		}
		c.JSON(http.StatusOK, resp)
	default:
		resp := responses.ErrorResponse{
			Error: "Invalid sign up type provided",
		}
		c.JSON(http.StatusBadRequest, resp)
	}
}

// CreateUser handles POST /users - create a new user along with (optional) phone numbers.
func CreateUser(c *gin.Context) {
	db := system.GetDb()
	var user model.User
	// Bind the incoming JSON payload to the user struct.
	if err := c.ShouldBindJSON(&user); err != nil {
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Check if a user with the same email already exists.
	var existingUser model.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		// Record found - email already exists.
		resp := responses.ErrorResponse{
			Error: "Email already exists",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	} else if err != gorm.ErrRecordNotFound {
		// Some other error occurred while querying.
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// Set timestamps for the new record.
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Create the user along with any associated phone numbers.
	if err := db.Create(&user).Error; err != nil {
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// RLP - API

	resp := responses.APIResponse{
		Message: "user created",
		Data:    user,
	}
	c.JSON(http.StatusCreated, resp)
}
