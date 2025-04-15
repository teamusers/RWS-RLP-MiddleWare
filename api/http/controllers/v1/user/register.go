package user

import (
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
	"github.com/google/uuid"
)

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
		if errors.Is(err, services.ErrRecordNotFound) {
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
		resp := responses.APIResponse{
			Message: "invalid json request body",
			Data:    model.User{},
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	//TO DO - If sign_up_type = TM: request TM info and validate

	//TO DO - Update RLP_ID generation logic
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
	var req requests.User
	req.ExternalID = user.ExternalID
	req.ExternalTYPE = user.ExternalTYPE // Adjust if field names differ between the structs
	req.Email = user.Email
	req.BurnPin = user.BurnPin
	req.GR_ID = "gr_id"                         // To be update by rlp.gr_id
	req.RLP_ID = rlpId.String()                 // To be update by rlp.rlp_id
	req.RWS_Membership_ID = "rws_membership_id" // To be update by rws_membership_id
	req.RWS_Membership_Number = 123456          // To be update by RWS_Membership_Number

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
