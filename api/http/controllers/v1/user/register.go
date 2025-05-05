package user

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	var req requests.RegisterUser
	// Bind the incoming JSON payload to the user struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("BindJSON error:", err)
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	switch req.SignUpType {
	case "NEW": // No action for now
	case "GR-CMS":
		cachedProfile, err := system.ObjectGet(strconv.Itoa(req.RegId), &model.GrProfile{})
		if err != nil {
			log.Printf("error getting cache value: %v", err)
			c.JSON(http.StatusConflict, responses.ApiResponse[any]{
				Code:    codes.CACHED_PROFILE_NOT_FOUND,
				Message: "cached profile not found",
			})
			return
		}

		req.User = cachedProfile.MapGrProfileToLbeUser()

		// match tier (assuming "Class X" format)
		assignTier(&req.User, cachedProfile.Class, c)

	case "GR":
		req.User = req.GrProfile.MapGrProfileToLbeUser()

		// match tier (assuming "Class X" format)
		assignTier(&req.User, req.GrProfile.Class, c)

	case "TM":
		// TODO: Request and Validate TM info

		// match tier
		req.User.Tier = "Tier M"

	default:
		c.JSON(http.StatusBadRequest, responses.ApiResponse[any]{
			Code:    codes.INVALID_REQUEST_BODY,
			Message: "invalid sign_up_type",
		})
		return
	}

	newRlpNumbering, newRlpNumberingErr := utils.GenerateNextRLPUserNumberingWithRetry()
	if newRlpNumberingErr != nil {
		log.Printf("Generate RLP User Number failed: %v", newRlpNumberingErr)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	log.Printf("RLP User Number generated: %v", newRlpNumbering)

	//To DO - RLP : Test Actual RLP End Points
	profileResp, err := services.Profile("", req.User.MapLbeToRlpUser(newRlpNumbering.RLP_ID), "PUT", services.ProfileURL)
	if err != nil {
		// Log the error
		log.Printf("RLP Register User failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	//TO DO - RLP: Request User Tier update
	if req.User.Tier != "" {
		// call rlp
	}

	// Request member service update
	var memberReq requests.CreateMemberUser
	memberReq.User.Email = req.User.Email
	memberReq.User.RLP_ID = newRlpNumbering.RLP_ID
	memberReq.User.RLP_NO = newRlpNumbering.RLP_NO
	memberReq.User.GR_ID = req.GrProfile.Id

	errRegister := services.PostRegisterUser(memberReq)
	if errRegister != nil {
		// Log the error
		log.Printf("Post Register User failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	resp := responses.ApiResponse[responses.CreateUserResponseData]{
		Code:    codes.SUCCESSFUL,
		Message: "user created",
		Data: responses.CreateUserResponseData{
			User: profileResp.User.MapRlpToLbeUser(),
		},
	}
	c.JSON(http.StatusCreated, resp)

	// purge regId cache if used
	if req.SignUpType == "GR-CMS" { //TODO: Make into enum
		system.ObjectDelete(strconv.Itoa(req.RegId))
	}
}

// VerifyGrExistence godoc
// @Summary      Verify GR member existence
// @Description  Checks if a GR member ID is already registered.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      requests.VerifyGrUser   true  "GR registration check payload"
// @Success      200      {object}  responses.GrExistenceSuccessResponse  "GR member found"
// @Failure      400      {object}  responses.ErrorResponse                     "Invalid JSON request body"
// @Failure      401      {object}  responses.ErrorResponse                                       "Unauthorized – API key missing or invalid"
// @Failure      500      {object}  responses.ErrorResponse                            "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/gr [post]
func VerifyGrExistence(c *gin.Context) {
	var req requests.VerifyGrUser

	// Bind the incoming JSON payload.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	// TO DO - verifyMemberExistence by gr_id and return error if found (GR_MEMBER_LINKED)

	cmsMember, err := services.GRMemberProfile(req.GrProfile.Id, nil, "GET", services.GetMemberURL)
	if err != nil {
		// Log the error
		log.Printf("Error while getting GR Member: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	// return response from CMS
	resp := responses.ApiResponse[responses.VerifyGrUserResponseData]{
		Code:    codes.SUCCESSFUL,
		Message: "gr profile found",
		Data:    responses.VerifyGrUserResponseData{GrProfile: cmsMember.MapCmsToLbeGrProfile(), Otp: model.Otp{}},
	}
	c.JSON(http.StatusOK, resp)
}

// VerifyGrCmsExistence godoc
// @Summary      Verify and cache GR CMS member
// @Description  Checks if a GR CMS member email is in the system and caches their profile for follow‑up.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      requests.VerifyGrCmsUser  true  "GR CMS register payload"
// @Success      200      {object}  responses.GrCmsExistenceSuccessResponse{}          "Email not found; profile cached"
// @Failure      400      {object}  responses.ErrorResponse  "Invalid JSON request body"
// @Failure      401      {object}  responses.ErrorResponse                      "Unauthorized – API key missing or invalid"
// @Failure      409      {object}  responses.ErrorResponse                      "Email already registered"
// @Failure      500      {object}  responses.ErrorResponse               "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/gr-cms [post]
func VerifyGrCmsExistence(c *gin.Context) {
	var req requests.VerifyGrCmsUser

	// Bind the incoming JSON payload.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	respData, err := services.VerifyMemberExistence(req.GrProfile.Email, false)
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
		system.ObjectSet(regId.String(), req.GrProfile, 30*time.Minute)

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

	cachedGrCmsProfile, err := system.ObjectGet(regId, &model.GrProfile{})
	if err != nil {
		log.Printf("error getting cache value: %v", err)
		resp := responses.ApiResponse[any]{
			Code:    codes.CACHED_PROFILE_NOT_FOUND,
			Message: "cached profile not found",
		}
		c.JSON(http.StatusConflict, resp)
		return
	}

	// return dob from cached profile
	resp := responses.ApiResponse[responses.VerifyGrCmsUserResponseData]{
		Code:    codes.SUCCESSFUL,
		Message: "cached profile found",
		Data: responses.VerifyGrCmsUserResponseData{
			RegId: regId,
			GrProfile: model.GrProfile{
				DateOfBirth: cachedGrCmsProfile.DateOfBirth,
			},
		},
	}
	c.JSON(http.StatusOK, resp)
}

func assignTier(user *model.User, class string, c *gin.Context) {
	tier, err := GrTierMatching(class)
	if err != nil {
		log.Printf("error matching gr class to member tier: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}
	user.Tier = tier
}

func GrTierMatching(grClass string) (string, error) {
	parts := strings.Fields(grClass) // splits by whitespace
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid gr class format")
	}

	classLevel, _ := strconv.Atoi(parts[1])

	if classLevel < 1 {
		return "", fmt.Errorf("invalid gr class format")
	}

	if classLevel == 1 {
		return "", nil // if class level 1, return empty for Tier A
	} else if classLevel == 2 {
		return "Tier B", nil
	} else if classLevel >= 3 && classLevel <= 5 {
		return "Tier C", nil
	} else {
		return "Tier D", nil
	}
}
