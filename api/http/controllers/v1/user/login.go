package home

import (
	"context"
	"errors"
	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/codes"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUsers handles GET /users
func Login(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		resp := responses.ErrorResponse{
			Error: "Valid email is required as query parameter",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	user, err := services.GetLoginUserByEmail(email)
	if err != nil {
		if errors.Is(err, services.ErrRecordNotFound) {
			// If no user is found, return an error.
			resp := responses.APIResponse{
				Message: "email not found",
				Data: responses.LoginResponse{
					OTP:               "",
					ExpireIn:          0,
					LoginSessionToken: "",
					LoginExpireIn:     0,
				},
			}
			c.JSON(codes.CODE_EMAIL_NOTFOUND, resp)
			return
		}
		// For any other errors, return an internal server error.
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// Assign the values from the service response.
	loginSessionToken := user.Data.LoginSessionToken
	loginExpireInSeconds := user.Data.LoginExpireIn

	// If the user exists, generate OTP using the service.
	otpService := services.NewOTPService()
	ctx := context.Background()
	otpResp, err := otpService.GenerateOTP(ctx, email)
	if err != nil {
		resp := responses.ErrorResponse{
			Error: "Failed to generate OTP",
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	//Call send email services
	resp := responses.APIResponse{
		Message: "email found",
		Data: responses.LoginResponse{
			OTP:               otpResp.OTP,
			ExpireIn:          otpResp.ExpiresAt,
			LoginSessionToken: loginSessionToken,
			LoginExpireIn:     loginExpireInSeconds,
		},
	}
	c.JSON(http.StatusOK, resp)
}
