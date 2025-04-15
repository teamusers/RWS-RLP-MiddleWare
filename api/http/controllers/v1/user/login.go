package user

import (
	"context"
	"errors"
	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/codes"
	"lbe/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

	var req requests.Login

	// Bind the incoming JSON payload to the req struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := responses.ErrorResponse{
			Error: "Invalid request payload",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Check if the email field is empty.
	if req.Email == "" {
		resp := responses.ErrorResponse{
			Error: "Valid email is required in the request body",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	user, err := services.GetLoginUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, services.ErrRecordNotFound) {
			resp := responses.APIResponse{
				Message: "email not found",
				Data: responses.LoginResponse{
					Otp: model.Otp{
						Otp:       nil,
						OtpExpiry: nil,
					},
					LoginSessionToken: model.LoginSessionToken{
						LoginSessionToken:       nil,
						LoginSessionTokenExpiry: nil,
					},
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

	otpService := services.NewOTPService()
	ctx := context.Background()
	otpResp, err := otpService.GenerateOTP(ctx, req.Email)
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
			Otp:               otpResp,
			LoginSessionToken: user.Data,
		},
	}
	c.JSON(http.StatusOK, resp)
}
