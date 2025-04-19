package model

// LoginSessionToken holds the session token and its expiry timestamp.
//
// Example:
// {
//   "login_session_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//   "login_session_token_expiry": 1744176000
// }
type LoginSessionToken struct {
	// LoginSessionToken is the JWT issued after successful authentication.
	// example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	LoginSessionToken *string `json:"login_session_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`

	// LoginSessionTokenExpiry is the Unix timestamp (seconds since epoch) when the token expires.
	// example: 1744176000
	LoginSessionTokenExpiry *int64 `json:"login_session_token_expiry" example:"1744176000"`
}
