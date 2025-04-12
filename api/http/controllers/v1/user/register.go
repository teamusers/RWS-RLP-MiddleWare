package home

import (
	"context"
	"net/http"
	"time"

	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/codes"
	model "lbe/models"
	"lbe/system"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUsers handles GET /users
// If a user with the provided email already exists, it returns an error that the email already exists.
// If no user is found, it continues to generate an OTP.
func GetUser(c *gin.Context) {
	var req requests.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := responses.ErrorResponse{
			Error: "Valid email and sign_up_type are required in the request body",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	email := req.Email
	signUpType := req.SignUpType

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

/*
// UpdateUser handles PUT /users/:id - update an existing user and optionally update phone numbers.
func UpdateUser(c *gin.Context) {
	db := system.GetDb()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var user model.User
	if err := db.Preload("PhoneNumbers").First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var updatedData model.User
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update core user fields
	user.ExternalID = updatedData.ExternalID
	user.OptedIn = updatedData.OptedIn
	user.ExternalTYPE = updatedData.ExternalTYPE
	user.Email = updatedData.Email
	user.DOB = updatedData.DOB
	user.Country = updatedData.Country
	user.FirstName = updatedData.FirstName
	user.LastName = updatedData.LastName
	user.BurnPin = updatedData.BurnPin
	user.UpdatedAt = time.Now()

	// Optionally, update phone numbers.
	// Here we replace all existing phone numbers if new ones are provided.
	if len(updatedData.PhoneNumbers) > 0 {
		// Delete the current phone numbers for this user.
		if err := db.Where("user_id = ?", user.ID).Delete(&model.UserPhoneNumber{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		user.PhoneNumbers = updatedData.PhoneNumbers
		// Optionally, iterate over the new phone numbers to set their timestamps and foreign key.
		now := time.Now()
		for idx := range user.PhoneNumbers {
			user.PhoneNumbers[idx].UserID = user.ID
			user.PhoneNumbers[idx].CreatedAt = now
			user.PhoneNumbers[idx].UpdatedAt = now
		}
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE /users/:id - delete a user and cascade delete phone numbers.
func DeleteUser(c *gin.Context) {
	db := system.GetDb()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	// The foreign key constraint (with ON DELETE CASCADE) in the database will handle the deletion of associated phone numbers.
	if err := db.Delete(&model.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
*/
