package responses

import "lbe/model"

type MemberAuthResponse struct {
	Message string `json:"message"`
	Data    struct {
		AccessToken string `json:"access_token"`
	} `json:"data"`
}

// UserResponse represents the expected response from the users endpoint.
// You can expand the fields based on what the endpoint returns.
type UserResponse struct {
	Message string                  `json:"message"`
	Data    model.LoginSessionToken `json:"data"`
}
