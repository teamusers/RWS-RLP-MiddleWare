package user

import (
	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/codes"
	"lbe/model"
	"lbe/system"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMemberProfile(c *gin.Context) {
	external_id := c.Param("external_id")
	if external_id == "" {
		resp := responses.ErrorResponse{
			Error: "Valid external_id is required as query parameter",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	//To DO - RLP : To be change to RLP view user. RLP - API, Temporary get from DB 1st
	db := system.GetDb()
	var user model.User
	err := db.Preload("PhoneNumbers").Where("external_id = ?", external_id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {

			resp := responses.APIResponse{
				Message: "member not found",
				Data:    nil,
			}
			c.JSON(codes.CODE_EMAIL_NOTFOUND, resp)
			return
		}
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := responses.APIResponse{
		Message: "successful",
		Data:    user,
	}
	c.JSON(http.StatusCreated, resp)
}

func UpdateMemberProfile(c *gin.Context) {
	external_id := c.Param("external_id")
	if external_id == "" {
		resp := responses.ErrorResponse{
			Error: "Valid external_id is required as query parameter",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	//To DO - RLP : To be change to RLP update user. RLP - API, Temporary update DB 1st
	db := system.GetDb()
	var user model.User
	err := db.Preload("PhoneNumbers").Where("external_id = ?", external_id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp := responses.APIResponse{
				Message: "member not found",
				Data:    nil,
			}
			c.JSON(codes.CODE_EMAIL_NOTFOUND, resp)
			return
		}
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// Bind the update input from the JSON body.
	// You could use a dedicated struct for the allowed update fields.
	// Here we're reusing model.User for simplicity.
	var updateData model.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		resp := responses.ErrorResponse{
			Error: "Invalid input: " + err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Update the user's profile using the received values.
	// The Updates method will perform a non-zero update on the provided fields.
	err = db.Model(&user).Updates(updateData).Error
	if err != nil {
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// Optional: Reload the user record with the associated phone numbers.
	err = db.Preload("PhoneNumbers").Where("external_id = ?", external_id).First(&user).Error
	if err != nil {
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := responses.APIResponse{
		Message: "update successfully",
		Data:    user,
	}
	c.JSON(http.StatusOK, resp)
}

func UpdateBurnPin(c *gin.Context) {
	var req requests.UpdateBurnPin

	// Bind the incoming JSON payload to the req struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := responses.APIResponse{
			Message: "invalid json request body",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err := services.UpdateBurnPin(req)
	if err != nil {
		log.Printf("error encountered updating burn pin: %v", err)
		resp := responses.APIResponse{
			Message: "update unsuccessful",
		}
		c.JSON(http.StatusCreated, resp)
		return
	}

	resp := responses.APIResponse{
		Message: "update successful",
	}
	c.JSON(http.StatusOK, resp)
}
