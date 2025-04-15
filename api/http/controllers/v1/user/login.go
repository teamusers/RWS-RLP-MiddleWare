package user

import (
	"errors"
	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/codes"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		if errors.Is(err, services.ErrRecordNotFound) {
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
