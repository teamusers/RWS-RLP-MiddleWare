package responses

type LoginResponse struct {
	OTP               string `json:"otp"`
	ExpireIn          int64  `json:"expireIn"`
	LoginSessionToken string `json:"loginSessionToken"`
	LoginExpireIn     int64  `json:"login_expireIn"`
}
