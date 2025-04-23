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

// GetUserProfile godoc
// @Summary      Get user profile
// @Description  Retrieves the profile (including phone numbers) for a given user by external_id.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        external_id  path      string                      true  "user external ID"
// @Success      200          {object}  responses.GetUserSuccessResponse       "user found"
// @Failure      400          {object}  responses.ErrorResponse  "Invalid or missing external_id path parameter"
// @Failure      401          {object}  responses.ErrorResponse                          "Unauthorized – API key missing or invalid"
// @Failure      409          {object}  responses.ErrorResponse                       "existing user not found"
// @Failure      500          {object}  responses.ErrorResponse              "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/{external_id} [get]
func GetUserProfile(c *gin.Context) {
	external_id := c.Param("external_id")
	if external_id == "" {
		c.JSON(http.StatusBadRequest, responses.InvalidQueryParametersErrorResponse())
		return
	}

	//To DO - RLP : To be change to RLP view user. RLP - API, Temporary get from DB 1st
	//memberResp, err := services.Member(external_id, nil, "GET")

	db := system.GetDb()
	var user model.User
	err := db.Preload("PhoneNumbers").Where("external_id = ?", external_id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusConflict, responses.DefaultResponse(codes.EXISTING_USER_NOT_FOUND, "existing user not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	resp := responses.ApiResponse[model.User]{
		Code:    codes.SUCCESSFUL,
		Message: "user found",
		Data:    user,
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateUserProfile godoc
// @Summary      Update user profile
// @Description  Updates a user's profile fields (non‐zero values in the JSON body).
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        external_id  path      string                      true  "user external ID"
// @Param        request      body      requests.User               true  "Profile fields to update"
// @Success      200          {object}  responses.UpdateUserSuccessResponse      "Update successful"
// @Failure      400          {object}  responses.ErrorResponse    "Invalid JSON request body"
// @Failure      401          {object}  responses.ErrorResponse                         "Unauthorized – API key missing or invalid"
// @Failure      409          {object}  responses.ErrorResponse                          "existing user not found"
// @Failure      500          {object}  responses.ErrorResponse                "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/{external_id} [put]
func UpdateUserProfile(c *gin.Context) {
	external_id := c.Param("external_id")
	if external_id == "" {
		c.JSON(http.StatusBadRequest, responses.InvalidQueryParametersErrorResponse())
		return
	}

	//To DO - RLP : To be change to RLP update user. RLP - API, Temporary update DB 1st
	//memberResp, err := services.Member(external_id, nil, "PUT")

	db := system.GetDb()
	var user model.User
	err := db.Preload("PhoneNumbers").Where("external_id = ?", external_id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusConflict, responses.DefaultResponse(codes.EXISTING_USER_NOT_FOUND, "existing user not found"))
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
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request      body      requests.UpdateBurnPin  true  "Email + new burn PIN"
// @Success      200          {object}  responses.UpdateUserSuccessResponse                      "update successful"
// @Failure      400          {object}  responses.ErrorResponse    "Invalid JSON request body"
// @Failure      401          {object}  responses.ErrorResponse                       "Unauthorized – API key missing or invalid"
// @Failure      500          {object}  responses.ErrorResponse             "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/pin [put]
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
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	c.JSON(http.StatusOK, responses.DefaultResponse(codes.SUCCESSFUL, "update successful"))
}
