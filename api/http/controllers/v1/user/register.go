package user

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/codes"
	"lbe/model"
	"lbe/system"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	email := c.Param("email")
	signUpType := c.Param("sign_up_type")

	if email == "" || signUpType == "" {
		resp := responses.ErrorResponse{
			Error: "Valid email and sign_up_type are required as query parameters",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err := services.GetRegisterUserByEmail(email, signUpType)
	if err != nil {
		if errors.Is(err, services.ErrRecordNotFound) {
			resp := responses.APIResponse{
				Message: "email registered",
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
			c.JSON(codes.CODE_EMAIL_REGISTERED, resp)
			return
		}
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// Generate OTP using the service.
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

	//To DO : Cache GR member info within expiry timestamp & Generate reg_ID

	//To DO : send email
	resp := responses.APIResponse{
		Message: "email not registered",
		Data:    otpResp,
	}
	c.JSON(http.StatusOK, resp)

}
func CreateUser(c *gin.Context) {
	db := system.GetDb()
	var user model.User
	// Bind the incoming JSON payload to the user struct.
	if err := c.ShouldBindJSON(&user); err != nil {
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	//To DO : To be change to RLP create user. RLP - API, Temporary Store into DB 1st
	if err := db.Create(&user).Error; err != nil {
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	//To DO : get RLP information and link accordingly
	var req requests.User
	req.ExternalID = user.ExternalID
	req.ExternalTYPE = user.ExternalTYPE // Adjust if field names differ between the structs
	req.Email = user.Email
	req.BurnPin = user.BurnPin
	req.GR_ID = "gr_id"                         // To be update by rlp.gr_id
	req.RLP_ID = "rlp_id"                       // To be update by rlp.rlp_id
	req.RWS_Membership_ID = "rws_membership_id" // To be update by rws_membership_id
	req.RWS_Membership_Number = 123456          // To be update by RWS_Membership_Number

	err := services.PostRegisterUser(req)
	if err != nil {
		// Log the error
		log.Printf("Post Register User failed: %v", err)
		resp := responses.ErrorResponse{
			Error: err.Error(),
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
