package user

import (
	"context"
	"errors"
	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/codes"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	loginSessionToken := user.Data.LoginSessionToken
	loginExpireInSeconds := user.Data.LoginExpireIn

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

	//To DO : Call send email services

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
