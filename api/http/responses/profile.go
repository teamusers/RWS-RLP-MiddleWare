package responses

import "lbe/model"

type GetUserProfileResponse struct {
	User model.User `json:"user"`
}
