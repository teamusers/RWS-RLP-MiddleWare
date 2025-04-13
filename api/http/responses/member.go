package responses

type MemberAuthResponse struct {
	Message string `json:"message"`
	Data    struct {
		AccessToken string `json:"accessToken"`
	} `json:"data"`
}

// LoginData represents the token details returned when the email is found.
type LoginData struct {
	LoginSessionToken string `json:"loginSessionToken"`
	LoginExpireIn     int64  `json:"login_expireIn"`
}

// UserResponse represents the expected response from the users endpoint.
// You can expand the fields based on what the endpoint returns.
type UserResponse struct {
	Message string    `json:"message"`
	Data    LoginData `json:"data"`
}
