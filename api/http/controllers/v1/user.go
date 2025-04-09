package home

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rlp-member-service/api/http/services"
	model "rlp-member-service/models"
	"rlp-member-service/system"
)

// GetUsers handles GET /users - list all users along with their phone numbers.
/**
func GetUsers(c *gin.Context) {
	db := system.GetDb()
	var users []model.User
	// Preload phone numbers to include them in the JSON response
	if err := db.Preload("PhoneNumbers").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
**/

// SignUpRequest represents the expected JSON structure for the request body.
type SignUpRequest struct {
	Email      string `json:"email" binding:"required,email"`
	SignUpType string `json:"sign_up_type" binding:"required"`
}

// GetUsers handles GET /users
// If a user with the provided email already exists, it returns an error that the email already exists.
// If no user is found, it continues to generate an OTP.
func GetUser(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid email and sign_up_type are required in the request body"})
		return
	}
	email := req.Email
	signUpType := req.SignUpType

	switch signUpType {
	case "new":
		// Get a database handle.
		db := system.GetDb()
		// Attempt to find a user by email.
		var user model.User
		err := db.Preload("PhoneNumbers").Where("email = ?", email).First(&user).Error
		if err == nil {
			// User found: return an error indicating that the email already exists.
			c.JSON(201, gin.H{
				"message": "email registered",
				"data": gin.H{
					"otp":          nil,
					"otp_expireIn": nil,
				},
			})
			return
		}
		if err != gorm.ErrRecordNotFound {
			// An unexpected error occurred during the query.
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// If we reach here, it means user was not found.
		// Generate OTP using the service.
		otpService := services.NewOTPService()
		ctx := context.Background()
		otpResp, err := otpService.GenerateOTP(ctx, email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
			return
		}

		//send email

		// Return the response with the custom JSON format.
		c.JSON(http.StatusOK, gin.H{
			"message": "email not registered",
			"data": gin.H{
				"otp":          otpResp.OTP,
				"otp_expireIn": otpResp.ExpiresAt,
			},
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sign up type provided"})
	}
}

// CreateUser handles POST /users - create a new user along with (optional) phone numbers.
func CreateUser(c *gin.Context) {
	db := system.GetDb()
	var user model.User
	// Bind the incoming JSON payload to the user struct.
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a user with the same email already exists.
	var existingUser model.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		// Record found - email already exists.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		// Some other error occurred while querying.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set timestamps for the new record.
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Create the user along with any associated phone numbers.
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// RLP - API

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"data":    user,
	})
}

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
