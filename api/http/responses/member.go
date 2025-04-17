package responses

import "lbe/model"

type MemberAuthResponse struct {
	Message string `json:"message"`
	Data    struct {
		AccessToken string `json:"access_token"`
	} `json:"data"`
}

type MemberLoginResponse struct {
	Message string                  `json:"message"`
	Data    model.LoginSessionToken `json:"data"`
}

type GetMemberUserResponse struct {
	Message string                  `json:"message"`
	Data    model.MembershipUser 		`json:"data"`
}

type UpdateBurnPin struct {
	Message string                  `json:"message"`
}
