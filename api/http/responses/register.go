package responses

import (
	"lbe/model"
)

// GetGrMemberResponse is returned after verifying a GR member’s existence and issuing an OTP.
//
// Example:
//
//	{
//	  "gr_id": "GR12345",
//	  "f_name": "Jane",
//	  "l_name": "Doe",
//	  "email": "jane.doe@example.com",
//	  "dob": "1985-04-12",
//	  "mobile": "98765432",
//	  "otp": "654321",
//	  "otp_expiry": 1744176000
//	}
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
