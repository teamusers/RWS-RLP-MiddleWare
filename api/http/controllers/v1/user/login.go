package user

import (
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
// @Success      200      {object}  responses.LoginSuccessResponse  "Email found; OTP generated and sent; login session token returned"
// @Failure      400      {object}  responses.ErrorResponse            "Invalid JSON request body"
// @Failure      401      {object}  responses.ErrorResponse                             "Unauthorized – API key missing or invalid"
// @Failure      409      {object}  responses.ErrorResponse                             "existing user not found"
// @Failure      500      {object}  responses.ErrorResponse                     "Internal server error"
// @Security     ApiKeyAuth
// @Router       /user/login [post]
func Login(c *gin.Context) {

	var req requests.Login

	// Bind the incoming JSON payload to the req struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.InvalidRequestBodyErrorResponse())
		return
	}

	respData, err := services.VerifyMemberExistence(req.Email, true)
	if err != nil {
		log.Printf("error encountered verifying user existence: %v", err)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}

	switch respData.Code {
	case codes.FOUND:
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

		resp := responses.ApiResponse[responses.LoginResponseData]{
			Code:    codes.SUCCESSFUL,
			Message: "login successful",
			Data: responses.LoginResponseData{
				Otp:               otpResp,
				LoginSessionToken: respData.Data,
			},
		}
		c.JSON(http.StatusOK, resp)
		return

	case codes.NOT_FOUND:
		c.JSON(http.StatusConflict, responses.ExistingUserNotFoundErrorResponse())
		return

	default:
		log.Printf("error encountered getting login user: %v", respData.Message)
		c.JSON(http.StatusInternalServerError, responses.InternalErrorResponse())
		return
	}
}
