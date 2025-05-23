package requests

import model "lbe/model"

// VerifyUseExistenceRequest is the payload to verify if an email is already registered.
// If not registered, an OTP will be sent to this email.
type VerifyUserExistence struct {
	// Email address to check for existing registration.
	Email string `json:"email" binding:"required" example:"user@example.com"`
}

type RegisterUser struct {
	User       model.User `json:"user"`
	SignUpType string     `json:"sign_up_type" example:"NEW"`
	RegId      int        `json:"reg_id" example:"123456"`
}

type VerifyGrUser struct {
	User model.User `json:"user" binding:"required"`
}

type VerifyGrCmsUser struct {
	User model.User `json:"user" binding:"required"`
}
