package responses

import (
	"lbe/model"
)

const (
	// error codes
	RlpErrorCodeUserNotFound = "user_not_found"
)

// GetUserResponse represents the top-level JSON
type GetUserResponse struct {
	Status string            `json:"status"`
	User   model.RlpUserResp `json:"user"`
}

type UserProfileErrorResponse struct {
	Status string    `json:"status"`
	Errors RlpErrors `json:"errors"`
}

type RlpErrors struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
