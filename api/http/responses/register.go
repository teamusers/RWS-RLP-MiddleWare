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
type GetGrMemberResponseData struct {
	// GrMember contains the GR profile fields.

	User *GRProfilePayload `json:"user"`
	// Otp contains the one‑time password details.
	model.Otp
}
