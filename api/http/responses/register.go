package responses

import (
	"lbe/model"
)

type VerifyGrUserResponseData struct {
	// User contains user data
	User model.User `json:"user"`
	// Otp contains the one‑time password details.
	model.Otp
}

type VerifyGrCmsUserResponseData struct {
	RegId       string     `json:"reg_id"`
	DateOfBirth model.Date `json:"dob"`
}

type CreateUserResponseData struct {
	User model.User `json:"user"`
}
