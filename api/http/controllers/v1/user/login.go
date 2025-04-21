package user

import (
	"errors"
	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/codes"
	"lbe/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary      Start login flow via email
// @Description  Validates user email, generates an OTP, emails it, and returns the OTP details plus a login session token.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request  body      requests.Login          true  "Login request payload"
// @Success      200      {object}  responses.APIResponse{data=responses.LoginResponse}
// @Failure      400      {object}  responses.APIResponse   "bad request"
// @Failure      401      {object}  responses.APIResponse	"unauthorized"
// @Failure      500      {object}  responses.APIResponse   "internal error"
// @Security     ApiKeyAuth
// @Router       /user/login [post]
func Login(c *gin.Context) {

	var req requests.Login

	// Bind the incoming JSON payload to the req struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := responses.APIResponse{
			Message: "invalid json request body",
			Data:    responses.LoginResponse{},
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	user, err := services.GetLoginUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, services.ErrCondition) {
			resp := responses.APIResponse{
				Message: "email not found",
				Data:    responses.LoginResponse{},
			}
			c.JSON(codes.CODE_EMAIL_NOTFOUND, resp)
			return
		}
		log.Printf("error encountered getting login user: %v", err)
		resp := responses.APIResponse{
			Message: "internal error",
			Data:    responses.LoginResponse{},
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	otpService := services.NewOTPService()
	otpResp, err := otpService.GenerateOTP(c, req.Email)
	if err != nil {
		log.Printf("error encountered generating otp: %v", err)
		resp := responses.APIResponse{
			Message: "internal error",
			Data:    responses.LoginResponse{},
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
			Data:    responses.LoginResponse{},
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := responses.APIResponse{
		Message: "email found",
		Data: responses.LoginResponse{
			Otp:               otpResp,
			LoginSessionToken: user.Data,
		},
	}
	c.JSON(http.StatusOK, resp)
}
