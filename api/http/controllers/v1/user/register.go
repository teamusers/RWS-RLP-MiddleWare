package user

import (
	"encoding/json"
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

const (
	// signUpType enums
	signUpTypeNew   = "NEW"
	signUpTypeGR    = "GR"
	signUpTypeGRCMS = "GR_CMS"
	signUpTypeTM    = "TM"
)

// VerifyUserExistence godoc
// @Summary      Verify email for registration
// @Description  Checks if an email is already registered; if not, sends an OTP for signup.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      requests.VerifyUserExistence true  "Registration request payload"
// @Success      200      {object}  responses.RegisterSuccessResponse "existing user not found"
// @Failure      400      {object}  responses.ErrorResponse  "Invalid JSON request body"
// @Failure      401      {object}  responses.ErrorResponse                       "Unauthorized – API key missing or invalid"
// @Failure      409      {object}  responses.ErrorResponse                      "existing user found"
// @Failure      500      {object}  responses.ErrorResponse               "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/register/verify [post]
func VerifyUserExistence(c *gin.Context) {
	httpClient := utils.GetHttpClient(c.Request.Context())
	var req requests.VerifyUserExistence

	// Bind the incoming JSON payload.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	if respData, _, err := services.GetCIAMUserByEmail(c, httpClient, req.Email); err != nil {
		log.Printf("error encountered verifying user existence: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	} else if len(respData.Value) != 0 {
		c.JSON(http.StatusConflict, responses.DefaultResponse(codes.EXISTING_USER_FOUND, "existing user found"))
		return
	}

	log.Printf("user %s not found, generating otp", req.Email)

	// if user is not found, generate OTP
	otpService := services.NewOTPService()
	otpResp, err := otpService.GenerateOTP(c, req.Email)
	if err != nil {
		log.Printf("error encountered generating otp: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	// Send OTP email via ACS
	acsRequest := requests.AcsSendEmailByTemplateRequest{
		Email:   req.Email,
		Subject: services.AcsEmailSubjectRequestOtp,
		Data: requests.RequestEmailOtpTemplateData{
			Email: req.Email,
			Otp:   *otpResp.Otp,
		},
	}

	if err := services.PostAcsSendEmailByTemplate(c, httpClient, services.AcsEmailTemplateRequestOtp, acsRequest); err != nil {
		log.Printf("failed to send email otp: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	resp := responses.ApiResponse[model.Otp]{
		Code:    codes.SUCCESSFUL,
		Message: "existing user not found",
		Data:    otpResp,
	}
	c.JSON(http.StatusOK, resp)
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
	httpClient := utils.GetHttpClient(c.Request.Context())
	var req requests.RegisterUser
	// Bind the incoming JSON payload to the user struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("BindJSON error:", err)
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	switch req.SignUpType {
	case signUpTypeNew:
		req.User.Tier = "Tier A" // set to base tier
	case signUpTypeGRCMS:
		cachedProfile, err := system.ObjectGet(strconv.Itoa(req.RegId), &model.User{})
		if err != nil {
			log.Printf("error getting cache value: %v", err)
			c.JSON(http.StatusConflict, responses.ApiResponse[any]{
				Code:    codes.CACHED_PROFILE_NOT_FOUND,
				Message: "cached profile not found",
			})
			return
		}

		req.User = *cachedProfile

		// match tier (assuming "Class X" format)
		assignTier(&req.User, c)

	case signUpTypeGR:
		// match tier (assuming "Class X" format)
		assignTier(&req.User, c)

	case signUpTypeTM:
		// TODO: Request and Validate TM info
		req.User.UserProfile.EmployeeNumber = "TBC"

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

	// populate registrations defaults
	req.User.PopulateIdentifiers(newRlpNumbering.RLP_ID, newRlpNumbering.RLP_NO)

	rlpCreateUserRequest := req.User.MapLbeToRlpUser()
	rlpCreateUserRequest.OptedIn = true
	rlpCreateUserRequest.ExternalID = newRlpNumbering.RLP_ID
	rlpCreateUserRequest.ExternalIDType = "rlp_id"
	rlpCreateUserRequest.UserProfile.LanguagePreference = "EN"
	rlpCreateUserRequest.UserProfile.PreviousEmail = rlpCreateUserRequest.Email

	//To DO - RLP : Test Actual RLP End Points
	profileResp, _, err := services.PutProfile(c, httpClient, "", rlpCreateUserRequest)
	if err != nil {
		// Log the error
		log.Printf("RLP Register User failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	// RLP: Request User Tier update
	// TODO: Update to actual spec
	log.Println("RLP Trigger Update User Tier Event")
	userTierReq := requests.UserTierUpdateEventRequest{
		EventLookup: services.RlpEventNameUpdateUserTier,
		UserId:      newRlpNumbering.RLP_ID,
		UserTier:    req.User.Tier,
	}

	if _, _, err := services.UpdateUserTier(c, httpClient, userTierReq); err != nil {
		log.Printf("RLP Update User Tier failed: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	} else {
		profileResp.User.Tier = req.User.Tier // update tier for response dto
	}

	// Create CIAM User
	if respData, raw, err := services.PostCIAMRegisterUser(c, httpClient, requests.GenerateInitialRegistrationRequest(&req.User)); err != nil {
		// Log the error
		log.Printf("CIAM Register User failed: %v", err)

		var errResp responses.GraphApiErrorResponse
		if err := json.Unmarshal(raw, &errResp); err == nil {
			if errResp.Error.Message == responses.CiamUserAlreadyExists {
				c.JSON(http.StatusConflict, responses.ExistingUserFoundErrorResponse())
				return
			}
		}
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	} else {
		// add schema extensions
		grID := ""
		if req.User.GrProfile != nil {
			grID = req.User.GrProfile.Id
		}

		schemaExtensionsPayload := map[string]any{
			config.GetConfig().Api.Eeid.UserIdLinkExtensionKey: requests.UserIdLinkSchemaExtensionFields{
				RlpId: newRlpNumbering.RLP_ID,
				RlpNo: newRlpNumbering.RLP_NO,
				GrId:  grID,
			},
		}

		if _, err := services.PatchCIAMAddUserSchemaExtensions(c, httpClient, respData.Id, schemaExtensionsPayload); err != nil {
			log.Printf("CIAM Patch User Schema Extensions failed: %v", err)
			c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
			return
		}
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
	if req.SignUpType == signUpTypeGRCMS {
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
// @Success      200      {object}  responses.GrExistenceSuccessResponse  "gr profile found"
// @Failure      400      {object}  responses.ErrorResponse                     "Invalid JSON request body"
// @Failure      401      {object}  responses.ErrorResponse                                       "Unauthorized – API key missing or invalid"
// @Failure      500      {object}  responses.ErrorResponse                            "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/gr [post]
func VerifyGrExistence(c *gin.Context) {
	httpClient := utils.GetHttpClient(c.Request.Context())
	var req requests.VerifyGrUser

	// Bind the incoming JSON payload.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodySpecificErrorResponse(err.Error()))
		return
	}

	// verify if gr ID is unused
	if respData, _, err := services.GetCIAMUserByGrId(c, httpClient, req.User.GrProfile.Id); err != nil {
		log.Printf("error encountered verifying user existence: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	} else if len(respData.Value) != 0 {
		c.JSON(http.StatusConflict, responses.GrMemberIdLinkedErrorResponse())
		return
	}

	//TODO: add conflict response if cms member not found
	cmsMember, err := services.GRMemberProfile(req.User.GrProfile.Id, nil, "GET", services.GetMemberURL)
	if err != nil {
		// Log the error
		log.Printf("Error while getting GR Member: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	// generate OTP
	otpService := services.NewOTPService()
	otpResp, err := otpService.GenerateOTP(c, cmsMember.EmailAddress)
	if err != nil {
		log.Printf("error encountered generating otp: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	// send email otp via acs if gr member email consent = true
	if cmsMember.ContactOptionEmail {
		acsRequest := requests.AcsSendEmailByTemplateRequest{
			Email:   cmsMember.EmailAddress,
			Subject: services.AcsEmailSubjectRequestOtp,
			Data: requests.RequestEmailOtpTemplateData{
				Email: cmsMember.EmailAddress,
				Otp:   *otpResp.Otp,
			},
		}

		httpClient := utils.GetHttpClient(c.Request.Context())
		if err := services.PostAcsSendEmailByTemplate(c, httpClient, services.AcsEmailTemplateRequestOtp, acsRequest); err != nil {
			log.Printf("failed to send email otp: %v", err)
			c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
			return
		}
	}

	// return response from CMS
	resp := responses.ApiResponse[responses.VerifyGrUserResponseData]{
		Code:    codes.SUCCESSFUL,
		Message: "gr profile found",
		Data:    responses.VerifyGrUserResponseData{User: cmsMember.MapCmsProfileToLbeUser(), Otp: otpResp},
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
// @Success      200      {object}  responses.GrCmsExistenceSuccessResponse{}          "existing user not found"
// @Failure      400      {object}  responses.ErrorResponse  "Invalid JSON request body"
// @Failure      401      {object}  responses.ErrorResponse                      "Unauthorized – API key missing or invalid"
// @Failure      409      {object}  responses.ErrorResponse                      "Email already registered"
// @Failure      500      {object}  responses.ErrorResponse               "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/gr-cms [post]
func VerifyGrCmsExistence(c *gin.Context) {
	httpClient := utils.GetHttpClient(c.Request.Context())
	var req requests.VerifyGrCmsUser

	// Bind the incoming JSON payload.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodySpecificErrorResponse(err.Error()))
		return
	}

	if respData, _, err := services.GetCIAMUserByEmail(c, httpClient, req.User.Email); err != nil {
		log.Printf("error encountered verifying user existence: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	} else if len(respData.Value) != 0 {
		c.JSON(http.StatusConflict, responses.ExistingUserFoundErrorResponse())
		return
	}

	// verify if gr ID is unused
	if respData, _, err := services.GetCIAMUserByGrId(c, httpClient, req.User.GrProfile.Id); err != nil {
		log.Printf("error encountered verifying user existence: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	} else if len(respData.Value) != 0 {
		c.JSON(http.StatusConflict, responses.GrMemberIdLinkedErrorResponse())
		return
	}

	// TODO - Generate reg_id and cache gr member info within expiry timestamp
	regId := uuid.New()
	system.ObjectSet(regId.String(), req.User, 30*time.Minute)

	// generate url

	registrationUrl := fmt.Sprintf("%s/%v/%s", config.GetConfig().Api.Acs.GrCmsRegistrationUrlHost, time.Now().Unix(), regId)

	// TODO - send registration email with url and reg_id via acs
	acsRequest := requests.AcsSendEmailByTemplateRequest{
		Email:   req.User.Email,
		Subject: services.AcsEmailSubjectRequestOtp, //TODO: update subject
		Data: requests.RequestEmailOtpTemplateData{ //TODO: update template data
			Email: req.User.Email,
			Otp:   registrationUrl,
		},
	}

	//TODO: update template
	if err := services.PostAcsSendEmailByTemplate(c, httpClient, services.AcsEmailTemplateRequestOtp, acsRequest); err != nil {
		log.Printf("failed to send registration url email: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	// return email existence status
	c.JSON(http.StatusOK, responses.DefaultResponse(codes.SUCCESSFUL, "existing user not found"))
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

	cachedUserProfile, err := system.ObjectGet(regId, &model.User{})
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
			RegId:       regId,
			DateOfBirth: *cachedUserProfile.DateOfBirth,
		},
	}
	c.JSON(http.StatusOK, resp)
}

func assignTier(user *model.User, c *gin.Context) {
	tier, err := GrTierMatching(user.GrProfile.Class)
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
