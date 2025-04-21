package user

import (
	"errors"
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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// VerifyUserExistence godoc
// @Summary      Verify email for registration
// @Description  Checks if an email is already registered; if not, sends an OTP for signup.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      requests.Register true  "Registration request payload"
// @Success      200      {object}  responses.APIResponse{data=model.Otp}    "email not registered, OTP sent"
// @Failure      400      {object}  responses.APIResponse                    "invalid JSON"
// @Failure      409      {object}  responses.APIResponse                    "email already registered"
// @Failure      500      {object}  responses.APIResponse                    "internal error"
// @Security     ApiKeyAuth
// @Router       /user/register/verify [post]
func VerifyUserExistence(c *gin.Context) {
	var req requests.Register

	// Bind the incoming JSON payload.
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := responses.APIResponse{
			Message: "invalid json request body",
			Data:    model.Otp{},
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err := services.GetRegisterUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, services.ErrCondition) {
			resp := responses.APIResponse{
				Message: "email registered",
				Data:    model.Otp{},
			}
			c.JSON(codes.CODE_EMAIL_REGISTERED, resp)
			return
		}
		log.Printf("error encountered getting registered user: %v", err)
		resp := responses.APIResponse{
			Message: "internal error",
			Data:    model.Otp{},
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// Generate OTP using the service.
	otpService := services.NewOTPService()
	otpResp, err := otpService.GenerateOTP(c, req.Email)
	if err != nil {
		log.Printf("error encountered generating otp: %v", err)
		resp := responses.APIResponse{
			Message: "internal error",
			Data:    model.Otp{},
		}
		c.JSON(http.StatusInternalServerError, resp)
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
		resp := responses.APIResponse{
			Message: "internal error",
			Data:    model.User{},
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := responses.APIResponse{
		Message: "email not registered",
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
// @Param        user     body      model.User        true  "User create payload"
// @Success      201      {object}  responses.APIResponse{data=model.User}  "user created"
// @Failure      400      {object}  responses.APIResponse                   "invalid JSON"
// @Failure      500      {object}  responses.APIResponse                   "internal error"
// @Security     ApiKeyAuth
// @Router       /user/register [post]
func CreateUser(c *gin.Context) {
	db := system.GetDb()
	var user model.User
	// Bind the incoming JSON payload to the user struct.
	if err := c.ShouldBindJSON(&user); err != nil {
		// Log the *actual* bind error so we can see what’s wrong
		log.Printf("BindJSON error: %v", err)

		resp := responses.APIResponse{
			Message: fmt.Sprintf("invalid json request body: %v", err),
			Data:    model.User{},
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	//TO DO - If sign_up_type = TM: request TM info and validate

	//TO DO - Update RLP_ID generation logic - yyyymmdd 0000-9999
	rlpId := uuid.New()

	//TO DO - Add member tier matching logic

	//To DO - RLP : To be change to RLP create user. RLP - API, Temporary Store into DB 1st
	if err := db.Create(&user).Error; err != nil {
		resp := responses.APIResponse{
			Message: "internal server error",
			Data:    model.User{},
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	//TO DO - If sign_up_type = GR or TM: Request user tier update (redundant?)

	//To DO - RLP | member service : get RLP information and link accordingly to member service
	var req requests.CreateUser
	req.User.ExternalID = user.ExternalID
	//req.User.ExternalTYPE = user.ExternalTYPE // Adjust if field names differ between the structs
	req.User.Email = user.Email
	//req.User.BurnPin = user.BurnPin
	req.User.GR_ID = "gr_id"                         // To be update by rlp.gr_id
	req.User.RLP_ID = rlpId.String()                 // To be update by rlp.rlp_id
	req.User.RWS_Membership_ID = "rws_membership_id" // To be update by rws_membership_id
	req.User.RWS_Membership_Number = 123456          // To be update by RWS_Membership_Number
	req.Email = user.Email                           // To be update by RWS_Membership_Number

	//TO DO - Request member service update - different based on sign_up_type
	err := services.PostRegisterUser(req)
	if err != nil {
		// Log the error
		log.Printf("Post Register User failed: %v", err)
		resp := responses.APIResponse{
			Message: "internal server error",
			Data:    model.User{},
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := responses.APIResponse{
		Message: "user created",
		Data:    user,
	}
	c.JSON(http.StatusCreated, resp)
}

// VerifyGrExistence godoc
// @Summary      Verify GR member existence
// @Description  Checks if a GR member ID is already registered.
// @Tags         gr
// @Accept       json
// @Produce      json
// @Param        request  body      requests.RegisterGr   true  "GR registration check payload"
// @Success      200      {object}  responses.APIResponse                     "email registered"
// @Failure      400      {object}  responses.APIResponse                     "invalid JSON"
// @Failure      500      {object}  responses.APIResponse                     "internal error"
// @Security     ApiKeyAuth
// @Router       /gr/verify [post]
func VerifyGrExistence(c *gin.Context) {
	var req requests.RegisterGr

	// Bind the incoming JSON payload.
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := responses.APIResponse{
			Message: "invalid json request body",
			Data:    responses.GetGrMemberResponse{},
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// TO DO - CMS: Request GR member info

	// return response from CMS
	resp := responses.APIResponse{
		Message: "successful",
		Data: responses.GetGrMemberResponse{
			GrMember: model.GrMember{
				GrId: &req.GrId,
			},
		},
	}
	c.JSON(http.StatusOK, resp)

}

// VerifyGrCmsExistence godoc
// @Summary      Verify and cache GR CMS member
// @Description  Checks if a GR CMS member email is in the system and caches their profile for follow‑up.
// @Tags         grcms
// @Accept       json
// @Produce      json
// @Param        request  body      requests.RegisterGrCms  true  "GR CMS register payload"
// @Success      200      {object}  responses.APIResponse                     "email not registered"
// @Failure      400      {object}  responses.APIResponse                     "invalid JSON"
// @Failure      409      {object}  responses.APIResponse                     "email already registered"
// @Failure      500      {object}  responses.APIResponse                     "internal error"
// @Security     ApiKeyAuth
// @Router       /grcms/verify [post]
func VerifyGrCmsExistence(c *gin.Context) {
	var req requests.RegisterGrCms

	// Bind the incoming JSON payload.
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := responses.APIResponse{
			Message: "invalid json request body",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err := services.GetRegisterUserByEmail(*req.GrMember.Email)
	if err != nil {
		if errors.Is(err, services.ErrCondition) { //TO DO - Fix and standardize error
			resp := responses.APIResponse{
				Message: "email registered",
			}
			c.JSON(codes.CODE_EMAIL_REGISTERED, resp)
			return
		}
		log.Printf("error encountered getting registered user: %v", err)
		resp := responses.APIResponse{
			Message: "internal error",
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// TO DO - cache gr member info within expiry timestamp and generate reg_id
	regId := uuid.New()
	system.ObjectSet(regId.String(), req.GrMember, 30*time.Minute)

	// TO DO - send registration email with url and reg_id

	// return email existence status
	resp := responses.APIResponse{
		Message: "email not registered",
	}
	c.JSON(http.StatusOK, resp)

}

// GetCachedGrCmsProfile godoc
// @Summary      Get cached GR CMS profile
// @Description  Retrieves a temporarily cached GR CMS profile by registration ID.
// @Tags         grcms
// @Accept       json
// @Produce      json
// @Param        reg_id   path      string                  true  "Registration ID"
// @Success      200      {object}  responses.APIResponse{data=model.GrMember}  "successful"
// @Failure      400      {object}  responses.APIResponse                    "missing reg_id"
// @Failure      404      {object}  responses.APIResponse                    "not found or expired"
// @Security     ApiKeyAuth
// @Router       /grcms/profile/{reg_id} [get]
func GetCachedGrCmsProfile(c *gin.Context) {

	regId := c.Param("reg_id")
	if regId == "" {
		resp := responses.APIResponse{
			Message: "missing reg_id query parameter",
			Data:    model.GrMember{},
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	cachedGrCmsProfile, err := system.ObjectGet(regId, &model.GrMember{})
	if err != nil {
		log.Printf("error getting cache value: %v", err)
		resp := responses.APIResponse{
			Message: "unsuccessful",
			Data:    model.GrMember{},
		}
		c.JSON(http.StatusCreated, resp)
		return
	}

	// return cached profile
	resp := responses.APIResponse{
		Message: "successful",
		Data:    cachedGrCmsProfile,
	}
	c.JSON(http.StatusOK, resp)

	// delete cached data since value found
	system.ObjectDelete(regId)
}
