package model

type LoginSessionToken struct {
	LoginSessionToken       *string `json:"login_session_token"`
	LoginSessionTokenExpiry *int64  `json:"login_session_token_expiry"`
}
