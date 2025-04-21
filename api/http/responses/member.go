package responses

import "lbe/model"

// MemberAuthResponse is returned after exchanging valid credentials for an access token.
//
// Example:
// {
//   "message": "token successfully generated",
//   "data": {
//     "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
//   }
// }
type MemberAuthResponse struct {
	// Message provides a human‑readable status.
	// example: "token successfully generated"
	Message string `json:"message" example:"token successfully generated"`

	// Data holds the authentication token payload.
	Data struct {
		// AccessToken is the JWT issued to the client for subsequent requests.
		// example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
		AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	} `json:"data"`
}

// MemberLoginResponse is returned after a user successfully logs in.
//
// Example:
// {
//   "message": "login successful",
//   "data": {
//     "login_session_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//     "login_session_token_expiry": 1744176000
//   }
// }
type MemberLoginResponse struct {
	// Message provides a human‑readable status.
	// example: "login successful"
	Message string `json:"message" example:"login successful"`

	// Data contains the session token and its expiry.
	Data model.LoginSessionToken `json:"data"`
}

// GetMemberUserResponse is returned when fetching a member’s full profile.
//
// Example:
// {
//   "message": "successful",
//   "data": {
//     "id": 123,
//     "external_id": "abc123",
//     "email": "user@example.com",
//     // ... other fields ...
//   }
// }
type GetMemberUserResponse struct {
	// Message provides a human‑readable status.
	// example: "successful"
	Message string `json:"message" example:"successful"`

	// Data contains the detailed user profile.
	Data model.MembershipUser `json:"data"`
}

// UpdateBurnPinResponse is returned after updating the user’s burn PIN.
//
// Example:
// {
//   "message": "update successful"
// }
type UpdateBurnPin struct {
	// Message indicates whether the update was successful.
	// example: "update successful"
	Message string `json:"message" example:"update successful"`
}

type GetRlpMemberUserResponse struct {
	// Message provides a human‑readable status.
	// example: "successful"
	Message string `json:"message" example:"successful"`

	// Data contains the detailed user profile.
	Data model.User `json:"data"`
}
