package responses

import "lbe/model"

// LoginResponse contains the OTP and session token details returned after a successful login.
//
// Example:
// {
//   "otp": "123456",
//   "otp_expiry": 1744076000,
//   "login_session_token": "eyJhbGciOiJIUzI1...",
//   "login_session_token_expiry": 1744176000
// }
type LoginResponse struct {
	model.Otp
	model.LoginSessionToken
}
