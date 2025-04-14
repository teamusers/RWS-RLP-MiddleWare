package responses

import "lbe/model"

type LoginResponse struct {
	model.Otp
	model.LoginSessionToken
}
