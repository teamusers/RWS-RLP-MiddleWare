package user

import (
	"fmt"
	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/codes"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

	//To DO - RLP : Test Actual RLP End Points
	profileResp, err := services.Profile(external_id, nil, "GET", services.ProfileURL)
	if err != nil {
		// Log the error
		log.Printf("GET User Profile failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	resp := responses.ApiResponse[*responses.GetUserResponse]{
		Code:    codes.SUCCESSFUL,
		Message: "user found",
		Data:    profileResp,
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
// @Param        request      body      requests.UserRequest               true  "Profile fields to update"
// @Success      200          {object}  responses.UpdateUserSuccessResponse      "Update successful"
// @Failure      400          {object}  responses.ErrorResponse    "Invalid JSON request body"
// @Failure      401          {object}  responses.ErrorResponse                         "Unauthorized – API key missing or invalid"
// @Failure      409          {object}  responses.ErrorResponse                          "existing user not found"
// @Failure      500          {object}  responses.ErrorResponse                "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/update/{external_id} [put]
func UpdateUserProfile(c *gin.Context) {

	var user requests.UserRequest
	// Bind the incoming JSON payload to the user struct.
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("BindJSON error:", err)
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	external_id := c.Param("external_id")
	if external_id == "" {
		c.JSON(http.StatusBadRequest, responses.InvalidQueryParametersErrorResponse())
		return
	}

	//To DO - RLP : To be change to RLP update user. RLP - API, Temporary update DB 1st
	//memberResp, err := services.Member(external_id, nil, "PUT")
	//To DO - RLP : Test Actual RLP End Points
	profileResp, err := services.Profile(external_id, user, "PUT", services.ProfileURL)
	if err != nil {
		// Log the error
		log.Printf("Update User Profile failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	resp := responses.ApiResponse[*responses.GetUserResponse]{
		Code:    codes.SUCCESSFUL,
		Message: "update successful",
		Data:    profileResp,
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
