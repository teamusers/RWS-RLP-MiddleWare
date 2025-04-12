package responses

type SignUpResponse struct {
	OTP      string `json:"otp"`
	ExpireIn int64  `json:"expireIn"`
}
