package responses

import "lbe/model"

type AuthSuccessResponse struct {
	// in: body
	Code    int64            `json:"code" example:"1000"`
	Message string           `json:"message" example:"token successfully generated"`
	Data    AuthResponseData `json:"data"`
}

type LoginSuccessResponse struct {
	// in: body
	Code    int64             `json:"code" example:"1000"`
	Message string            `json:"message" example:"email found"`
	Data    LoginResponseData `json:"data"`
}

type RegisterSuccessResponse struct {
	// in: body
	Code    int64     `json:"code" example:"1000"`
	Message string    `json:"message" example:"email not found"`
	Data    model.Otp `json:"data"`
}

type CreateSuccessResponse struct {
	// in: body
	Code    int64      `json:"code" example:"1000"`
	Message string     `json:"message" example:"user created"`
	Data    model.User `json:"data"`
}

type GrExistenceSuccessResponse struct {
	// in: body
	Code    int64          `json:"code" example:"1000"`
	Message string         `json:"message" example:"successful"`
	Data    model.GrMember `json:"data"`
}

type GrCmsExistenceSuccessResponse struct {
	// in: body
	Code    int64          `json:"code" example:"1003"`
	Message string         `json:"message" example:"email not found"`
	Data    model.GrMember `json:"data"`
}

type CachedGrCmsSuccessResponse struct {
	// in: body
	Code    int64          `json:"code" example:"1002"`
	Message string         `json:"message" example:"cached profile found"`
	Data    model.GrMember `json:"data"`
}

type GetMemberSuccessResponse struct {
	// in: body
	Code    int64      `json:"code" example:"1002"`
	Message string     `json:"message" example:"member found"`
	Data    model.User `json:"data"`
}

type UpdateMemberSuccessResponse struct {
	// in: body
	Code    int64      `json:"code" example:"1000"`
	Message string     `json:"message" example:"update successful"`
	Data    model.User `json:"data"`
}

// ErrorResponse is the standard envelope for error responses.
// swagger:response ErrorResponse
type ErrorResponse struct {
	// in: body

	// Code is your internal API status code, e.g. 1002
	Code int64 `json:"code"`
	// Message is a human‑readable description, e.g. "invalid json request body"
	Message string `json:"message"`
	Data    string `json:"data"`
}
