package model

// Otp holds a one‑time password and its expiry timestamp.
//
// Example:
// {
//   "otp": "123456",
//   "otp_expiry": 1744176000
// }
type Otp struct {
	// Otp is the one‑time password sent to the user.
	// example: "123456"
	Otp *string `json:"otp" example:"123456"`

	// OtpExpiry is the Unix timestamp (seconds since epoch) when the OTP expires.
	// example: 1744176000
	OtpExpiry *int64 `json:"otp_expiry" example:"1744176000"`
}
