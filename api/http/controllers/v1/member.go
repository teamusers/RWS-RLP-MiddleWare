package v1

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

// GetMemberProfile godoc
// @Summary      Get member profile
// @Description  Retrieves the profile (including phone numbers) for a given member by external_id.
// @Tags         member
// @Accept       json
// @Produce      json
// @Param        external_id  path      string                      true  "Member external ID"
// @Success      200          {object}  responses.APIResponse{data=model.User}  "OK"
// @Failure      400          {object}  responses.ErrorResponse              "bad request"
// @Failure      401      	  {object}  responses.APIResponse			  "unauthorized"
// @Failure      500          {object}  responses.ErrorResponse              "internal error"
// @Security     ApiKeyAuth
// @Router       /member/{external_id} [get]
func GetMemberProfile(c *gin.Context) {
	external_id := c.Param("external_id")
	if external_id == "" {
		c.JSON(http.StatusBadRequest, responses.InvalidQueryParametersErrorResponse())
		return
	}

	//To DO - RLP : To be change to RLP view user. RLP - API, Temporary get from DB 1st
	db := system.GetDb()
	var user model.User
	err := db.Preload("PhoneNumbers").Where("external_id = ?", external_id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusConflict, responses.DefaultResponse(codes.NOT_FOUND, "member not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	resp := responses.ApiResponse[model.User]{
		Code:    codes.FOUND,
		Message: "member found",
		Data:    user,
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateMemberProfile godoc
// @Summary      Update member profile
// @Description  Updates a member’s profile fields (non‐zero values in the JSON body).
// @Tags         member
// @Accept       json
// @Produce      json
// @Param        external_id  path      string                      true  "Member external ID"
// @Param        request      body      requests.User               true  "Profile fields to update"
// @Success      200          {object}  responses.APIResponse{data=model.User}  "update successful"
// @Failure      400          {object}  responses.ErrorResponse              "bad request"
// @Failure      401      	  {object}  responses.APIResponse			  "unauthorized"
// @Failure      500          {object}  responses.ErrorResponse              "internal error"
// @Security     ApiKeyAuth
// @Router       /member/{external_id} [put]
func UpdateMemberProfile(c *gin.Context) {
	external_id := c.Param("external_id")
	if external_id == "" {
		c.JSON(http.StatusBadRequest, responses.InvalidQueryParametersErrorResponse())
		return
	}

	//To DO - RLP : To be change to RLP update user. RLP - API, Temporary update DB 1st
	db := system.GetDb()
	var user model.User
	err := db.Preload("PhoneNumbers").Where("external_id = ?", external_id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusConflict, responses.DefaultResponse(codes.NOT_FOUND, "member not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	// Bind the update input from the JSON body.
	// You could use a dedicated struct for the allowed update fields.
	// Here we're reusing model.User for simplicity.
	var updateData model.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	// Update the user's profile using the received values.
	// The Updates method will perform a non-zero update on the provided fields.
	err = db.Model(&user).Updates(updateData).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	// Optional: Reload the user record with the associated phone numbers.
	err = db.Preload("PhoneNumbers").Where("external_id = ?", external_id).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	resp := responses.ApiResponse[model.User]{
		Code:    codes.SUCCESSFUL,
		Message: "update successful",
		Data:    user,
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateBurnPin godoc
// @Summary      Update user burn PIN
// @Description  Updates the burn PIN for a given email address.
// @Tags         member
// @Accept       json
// @Produce      json
// @Param        request      body      requests.UpdateBurnPin  true  "Email + new burn PIN"
// @Success      200          {object}  responses.APIResponse            "update successful"
// @Failure      400          {object}  responses.APIResponse            "bad request"
// @Failure      401      	  {object}  responses.APIResponse			  "unauthorized"
// @Failure      500          {object}  responses.APIResponse            "update unsuccessful"
// @Security     ApiKeyAuth
// @Router       /member/burn-pin [put]
func UpdateBurnPin(c *gin.Context) {
	var req requests.UpdateBurnPin

	// Bind the incoming JSON payload to the req struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	err := services.UpdateBurnPin(req)
	if err != nil {
		log.Printf("error encountered updating burn pin: %v", err)
		c.JSON(http.StatusCreated, responses.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, responses.DefaultResponse(codes.SUCCESSFUL, "update successful"))
}
