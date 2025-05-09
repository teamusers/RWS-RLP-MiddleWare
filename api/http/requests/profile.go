package requests

import "lbe/model"

type UpdateUserProfile struct {
	User model.User `json:"user"`
}
