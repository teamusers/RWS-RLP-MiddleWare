package responses

type OTPResponse struct {
	OTP       string `json:"otp_code"`
	ExpiresAt int64  `json:"expires_at"` // Unix timestamp
}
