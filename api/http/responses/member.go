package responses

import "lbe/model"

type MemberAuthResponseData struct {
	AccessToken string `json:"access_token"`
}

type GetRlpMemberUserResponse struct {
	// Message provides a human‑readable status.
	// example: "successful"
	Message string `json:"message" example:"successful"`

	// Data contains the detailed user profile.
	Data model.User `json:"data"`
}
