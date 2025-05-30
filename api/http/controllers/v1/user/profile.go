package user

import (
	"fmt"
	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/codes"
	"lbe/model"
	"lbe/utils"
	"log"
	"net/http"
	"time"

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
	httpClient := utils.GetHttpClient(c.Request.Context())
	external_id := c.Param("external_id")

	// TODO - RLP : Test Actual RLP End Points
	profileResp, _, err := services.GetProfile(c, httpClient, external_id)
	if err != nil {
		// Log the error
		log.Printf("GET User Profile failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	resp := responses.ApiResponse[*responses.GetUserProfileResponse]{
		Code:    codes.SUCCESSFUL,
		Message: "user found",
		Data: &responses.GetUserProfileResponse{
			User: profileResp.User.MapRlpToLbeUser(),
		},
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
// @Param        request      body      requests.UpdateUserProfile               true  "Profile fields to update"
// @Success      200          {object}  responses.UpdateUserSuccessResponse      "Update successful"
// @Failure      400          {object}  responses.ErrorResponse    "Invalid JSON request body"
// @Failure      401          {object}  responses.ErrorResponse                         "Unauthorized – API key missing or invalid"
// @Failure      409          {object}  responses.ErrorResponse                          "existing user not found"
// @Failure      500          {object}  responses.ErrorResponse                "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/update/{external_id} [put]
func UpdateUserProfile(c *gin.Context) {
	httpClient := utils.GetHttpClient(c.Request.Context())
	var req requests.UpdateUserProfile
	// Bind the incoming JSON payload to the user struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("BindJSON error:", err)
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	external_id := c.Param("external_id")

	// TODO - RLP : Test Actual RLP End Points
	rlpUpdateUserReq := requests.UserProfileRequest{
		User: req.User.MapLbeToRlpUser(),
	}
	profileResp, _, err := services.UpdateProfile(c, httpClient, external_id, rlpUpdateUserReq)
	if err != nil {
		// Log the error
		log.Printf("Update User Profile failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	resp := responses.ApiResponse[*responses.GetUserProfileResponse]{
		Code:    codes.SUCCESSFUL,
		Message: "update successful",
		Data: &responses.GetUserProfileResponse{
			User: profileResp.User.MapRlpToLbeUser(),
		},
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
// func UpdateBurnPin(c *gin.Context) {
// 	var req requests.UpdateBurnPin

// 	// Bind the incoming JSON payload to the req struct.
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
// 		return
// 	}

// 	err := services.UpdateBurnPin(req)
// 	if err != nil {
// 		log.Printf("error encountered updating burn pin: %v", err)
// 		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
// 		return
// 	}

// 	c.JSON(http.StatusOK, responses.DefaultResponse(codes.SUCCESSFUL, "update successful"))
// }

// WithdrawUser godoc
// @Summary      Update user profile
// @Description  Updates a user's profile fields (non‐zero values in the JSON body).
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        external_id  path      string                      true  "user external ID"
// @Param        request      body      requests.UpdateUserProfile               true  "Profile fields to update"
// @Success      200          {object}  responses.UpdateUserSuccessResponse      "Update successful"
// @Failure      400          {object}  responses.ErrorResponse    "Invalid JSON request body"
// @Failure      401          {object}  responses.ErrorResponse                         "Unauthorized – API key missing or invalid"
// @Failure      409          {object}  responses.ErrorResponse                          "existing user not found"
// @Failure      500          {object}  responses.ErrorResponse                "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/archive/{external_id} [put]
func WithdrawUserProfile(c *gin.Context) {
	//TODO: add rollback
	httpClient := utils.GetHttpClient(c.Request.Context())
	external_id := c.Param("external_id")

	// Retrieve user profile from RLP
	rlpResp, _, err := services.GetProfile(c, httpClient, external_id)
	if err != nil {
		// Log the error
		log.Printf("GET User Profile failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	// Retrieve CIAM id
	ciamUserId := ""

	if respData, _, err := services.GetCIAMUserByEmail(c, httpClient, rlpResp.User.Email); err != nil {
		log.Printf("error encountered verifying user existence: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	} else if len(respData.Value) == 0 {
		c.JSON(http.StatusConflict, responses.ExistingUserNotFoundErrorResponse())
		return
	} else {
		ciamUserId = respData.Value[0].ID
	}

	// Update user profile to withdraw status
	rlpUserProfile := rlpResp.User
	now := time.Now()
	timestamp := now.Format("060102150405") // yyMMddHHmmss

	rlpUpdateUserReq := requests.UserProfileRequest{
		User: model.RlpUserReq{
			Email: fmt.Sprintf("%s.delete_%v", rlpUserProfile.Email, timestamp),
			UserProfile: model.UserProfile{
				ActiveStatus: "0",
				MarketingPreference: model.MarketingPreference{
					Push:   model.BoolPtr(false),
					Email:  model.BoolPtr(false),
					Mobile: model.BoolPtr(false),
				},
			},
		},
	}

	profileResp, _, err := services.UpdateProfile(c, httpClient, external_id, rlpUpdateUserReq)
	if err != nil {
		// Log the error
		log.Printf("Update User Profile to withdraw failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	// Update CIAM user to set accountEnabled = false
	ciamPayload := requests.GraphDisableAccountRequest{
		AccountEnabled: false,
	}
	if _, err := services.PatchCIAMUpdateUser(c, httpClient, ciamUserId, ciamPayload); err != nil {
		// Log the error
		log.Printf("Update CIAM User AccountEnabled to false failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	//TODO: update API call to ACS to proper template
	acsRequest := requests.AcsSendEmailByTemplateRequest{
		Email:   rlpResp.User.Email,                 // original email
		Subject: services.AcsEmailSubjectRequestOtp, //TODO: update subject
		Data: requests.RequestEmailOtpTemplateData{ //TODO: update template data
			Email: rlpResp.User.Email,
			Otp:   "withdraw email",
		},
	}

	//TODO: update template
	if err := services.PostAcsSendEmailByTemplate(c, httpClient, services.AcsEmailTemplateRequestOtp, acsRequest); err != nil {
		log.Printf("failed to send withdrawal email: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	resp := responses.ApiResponse[*responses.GetUserProfileResponse]{
		Code:    codes.SUCCESSFUL,
		Message: "update successful",
		Data: &responses.GetUserProfileResponse{
			User: profileResp.User.MapRlpToLbeUser(),
		},
	}
	c.JSON(http.StatusOK, resp)

}
