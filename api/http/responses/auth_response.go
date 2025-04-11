package responses

type LoginResponse struct {
	OTP               *string `json:"otp"`
	ExpireIn          *int64  `json:"expireIn"`
	LoginSessionToken *string `json:"loginSessionToken"`
	LoginExpireIn     *string `json:"login_expireIn"`
}
