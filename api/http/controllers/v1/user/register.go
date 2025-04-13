package home

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
	model "lbe/models"
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
			// If no user is found, return an error.
			resp := responses.APIResponse{
				Message: "email registered",
				Data: responses.LoginResponse{
					OTP:               "",
					ExpireIn:          0,
					LoginSessionToken: "",
					LoginExpireIn:     0,
				},
			}
			c.JSON(codes.CODE_EMAIL_REGISTERED, resp)
			return
		}
		// For any other errors, return an internal server error.
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
		Data: responses.SignUpResponse{
			OTP:      otpResp.OTP,
			ExpireIn: otpResp.ExpiresAt,
		},
	}
	c.JSON(http.StatusOK, resp)

}

// CreateUser handles POST /users - create a new user along with (optional) phone numbers.
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

	// Set timestamps for the new record.
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	// To DO : To be change to RLP create user. RLP - API, Temporary Store into DB 1st
	if err := db.Create(&user).Error; err != nil {
		resp := responses.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

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
