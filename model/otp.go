package model

type Otp struct {
	Otp       *string `json:"otp"`
	OtpExpiry *int64  `json:"otp_expiry"` // Unix timestamp
}
