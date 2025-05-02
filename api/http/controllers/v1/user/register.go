package user

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/codes"
	"lbe/config"
	"lbe/model"
	"lbe/system"
	"lbe/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// VerifyUserExistence godoc
// @Summary      Verify email for registration
// @Description  Checks if an email is already registered; if not, sends an OTP for signup.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      requests.VerifyUserExistence true  "Registration request payload"
// @Success      200      {object}  responses.RegisterSuccessResponse "Email not registered; OTP sent"
// @Failure      400      {object}  responses.ErrorResponse  "Invalid JSON request body"
// @Failure      401      {object}  responses.ErrorResponse                       "Unauthorized – API key missing or invalid"
// @Failure      409      {object}  responses.ErrorResponse                      "Email already registered"
// @Failure      500      {object}  responses.ErrorResponse               "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/register/verify [post]
func VerifyUserExistence(c *gin.Context) {
	var req requests.VerifyUserExistence

	// Bind the incoming JSON payload.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	respData, err := services.VerifyMemberExistence(req.Email, false)
	if err != nil {
		log.Printf("error encountered verifying user existence: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	switch respData.Code {
	case codes.NOT_FOUND:
		otpService := services.NewOTPService()
		otpResp, err := otpService.GenerateOTP(c, req.Email)
		if err != nil {
			log.Printf("error encountered generating otp: %v", err)
			c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
			return
		}

		//Call send email services
		emailData := services.EmailOtpTemplateData{
			Email: req.Email,
			OTP:   *otpResp.Otp,
		}

		cfg := config.GetConfig()
		emailService := services.NewEmailService(&cfg.Smtp)
		if err := emailService.SendOtpEmail(req.Email, emailData); err != nil {
			log.Printf("failed to send email otp: %v", err)
			c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
			return
		}

		resp := responses.ApiResponse[model.Otp]{
			Code:    codes.SUCCESSFUL,
			Message: "existing user not found",
			Data: model.Otp{
				Otp: otpResp.Otp,
			},
		}
		c.JSON(http.StatusOK, resp)
		return

	case codes.FOUND:
		c.JSON(http.StatusConflict, responses.DefaultResponse(codes.EXISTING_USER_FOUND, "existing user found"))
		return

	default:
		log.Printf("error encountered getting registered user: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

}

// CreateUser godoc
// @Summary      Create new user
// @Description  Registers a new user record in the system.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user     body      requests.RegisterUser        true  "User create payload"
// @Success      201      {object}  responses.CreateSuccessResponse  "User created successfully"
// @Failure      400      {object}  responses.ErrorResponse  "Invalid JSON request body"
// @Failure      401      {object}  responses.ErrorResponse                      "Unauthorized – API key missing or invalid"
// @Failure      500      {object}  responses.ErrorResponse              "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/register [post]
func CreateUser(c *gin.Context) {
	var user requests.RegisterUser
	// Bind the incoming JSON payload to the user struct.
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("BindJSON error:", err)
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	//TO DO - If sign_up_type = TM: request TM info and validate

	newRlpNumbering, newRlpNumberingErr := utils.GenerateNextRLPUserNumberingWithRetry()
	if newRlpNumberingErr != nil {
		log.Printf("Post Register User failed: %v", newRlpNumberingErr)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	fmt.Println(newRlpNumbering)

	//TO DO - Add member tier matching logic

	//To DO - RLP : Test Actual RLP End Points
	profileResp, err := services.Profile("", user, "POST", services.UpdateProfileURL)
	if err != nil {
		// Log the error
		log.Printf("Post Register User failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}
	//TO DO - (If member tier > basic && sign_up_type = GR) or (sign_up_type == TM)L request user tier update (TM = tier M)

	//To DO - RLP | member service : get RLP information and link accordingly to member service
	var req requests.CreateUser
	//req.User.ExternalID = user.Users.ExternalID
	//req.User.ExternalTYPE = user.ExternalTYPE // Adjust if field names differ between the structs
	req.User.Email = user.Users.Email
	//req.User.BurnPin = user.BurnPin
	req.User.GR_ID = "gr_id"                 // To be update by rlp.gr_id
	req.User.RLP_ID = newRlpNumbering.RLP_ID // To be update by RLP_ID
	req.User.RLP_NO = newRlpNumbering.RLP_NO // To be rename & update by RLP_NO

	//TO DO - Request member service update - different based on sign_up_type
	errRegister := services.PostRegisterUser(req)
	if errRegister != nil {
		// Log the error
		log.Printf("Post Register User failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	resp := responses.ApiResponse[*responses.GetUserResponse]{
		Code:    codes.SUCCESSFUL,
		Message: "user created",
		Data:    profileResp,
	}
	c.JSON(http.StatusCreated, resp)
}

// VerifyGrExistence godoc
// @Summary      Verify GR member existence
// @Description  Checks if a GR member ID is already registered.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      requests.RegisterGr   true  "GR registration check payload"
// @Success      200      {object}  responses.GrExistenceSuccessResponse  "GR member found"
// @Failure      400      {object}  responses.ErrorResponse                     "Invalid JSON request body"
// @Failure      401      {object}  responses.ErrorResponse                                       "Unauthorized – API key missing or invalid"
// @Failure      500      {object}  responses.ErrorResponse                            "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/gr [post]
func VerifyGrExistence(c *gin.Context) {
	var req requests.RegisterGr

	// Bind the incoming JSON payload.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	// TO DO - verifyMemberExistence by gr_id and return error if found (GR_MEMBER_LINKED)

	cmsMember, err := services.GRMemberProfile(req.GrId, nil, "GET", services.GetMemberURL)
	if err != nil {
		// Log the error
		log.Printf("Error while getting GR Member: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	// return response from CMS
	// TODO: Fix
	resp := responses.ApiResponse[responses.GetGrMemberResponseData]{
		Code:    codes.SUCCESSFUL,
		Message: "gr profile found",
		Data:    responses.GetGrMemberResponseData{User: cmsMember, Otp: model.Otp{}},
	}
	c.JSON(http.StatusOK, resp)
}

// VerifyGrCmsExistence godoc
// @Summary      Verify and cache GR CMS member
// @Description  Checks if a GR CMS member email is in the system and caches their profile for follow‑up.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      requests.RegisterGrCms  true  "GR CMS register payload"
// @Success      200      {object}  responses.GrCmsExistenceSuccessResponse{}          "Email not found; profile cached"
// @Failure      400      {object}  responses.ErrorResponse  "Invalid JSON request body"
// @Failure      401      {object}  responses.ErrorResponse                      "Unauthorized – API key missing or invalid"
// @Failure      409      {object}  responses.ErrorResponse                      "Email already registered"
// @Failure      500      {object}  responses.ErrorResponse               "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/gr-cms [post]
func VerifyGrCmsExistence(c *gin.Context) {
	var req requests.RegisterGrCms

	// Bind the incoming JSON payload.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	respData, err := services.VerifyMemberExistence(*req.GrMember.Email, false)
	if err != nil {
		log.Printf("error encountered verifying user existence: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	switch respData.Code {
	case codes.NOT_FOUND:
		// TO DO - verifyMemberExistence by gr_id and return error if found (GR_MEMBER_LINKED)

		// TO DO - cache gr member info within expiry timestamp and generate reg_id
		regId := uuid.New()
		system.ObjectSet(regId.String(), req.GrMember, 30*time.Minute)

		// TO DO - send registration email with url and reg_id

		// return email existence status
		c.JSON(http.StatusOK, responses.DefaultResponse(codes.SUCCESSFUL, "existing user not found"))
		return

	case codes.FOUND:
		c.JSON(http.StatusConflict, responses.ExistingUserFoundErrorResponse())
		return

	default:
		log.Printf("error encountered getting registered user: %v", respData.Message)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}
}

// GetCachedGrCmsProfile godoc
// @Summary      Get cached GR CMS profile
// @Description  Retrieves a temporarily cached GR CMS profile by registration ID.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        reg_id   path      string                  true  "Registration ID"
// @Success      200      {object}  responses.CachedGrCmsSuccessResponse  "Cached profile found"
// @Failure      400      {object}  responses.ErrorResponse  "Registration ID is required"
// @Failure      409      {object}  responses.ErrorResponse                             "Cached profile not found"
// @Failure      500      {object}  responses.ErrorResponse                   "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/gr-reg/{reg_id} [get]
func GetCachedGrCmsProfile(c *gin.Context) {

	regId := c.Param("reg_id")
	if regId == "" {
		c.JSON(http.StatusBadRequest, responses.InvalidQueryParametersErrorResponse())
		return
	}

	cachedGrCmsProfile, err := system.ObjectGet(regId, &model.GrMember{})
	if err != nil {
		log.Printf("error getting cache value: %v", err)
		resp := responses.ApiResponse[any]{
			Code:    codes.CACHED_PROFILE_NOT_FOUND,
			Message: "cached profile not found",
		}
		c.JSON(http.StatusConflict, resp)
		return
	}

	// return cached profile
	resp := responses.ApiResponse[model.GrMember]{
		Code:    codes.SUCCESSFUL,
		Message: "cached profile found",
		Data:    *cachedGrCmsProfile,
	}
	c.JSON(http.StatusOK, resp)

	// delete cached data since value found
	system.ObjectDelete(regId)
}
