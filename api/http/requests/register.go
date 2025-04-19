package requests

import "lbe/model"

// RegisterRequest is the payload to verify if an email is already registered.
// If not registered, an OTP will be sent to this email.
type Register struct {
	// Email address to check for existing registration.
	Email string `json:"email" binding:"required" example:"user@example.com"`
}

// RegisterGrRequest is the payload to verify if a GR member ID is already registered.
type RegisterGr struct {
	// GR system identifier for the member.
	GrId string `json:"gr_id" binding:"required" example:"GR12345"`

	// PIN code associated with the GR member.
	GrPin string `json:"gr_pin" binding:"required" example:"9876"`
}

// RegisterGrCmsRequest is the payload to verify and cache a GR CMS member profile.
// Embeds the GrMember model for profile fields and includes a callback URL.
type RegisterGrCms struct {
	// GR member profile fields to cache.
	model.GrMember

	// URL to send the registration confirmation link to.
	Url string `json:"url" binding:"required" example:"https://example.com/confirm?reg_id=abc123"`
}
