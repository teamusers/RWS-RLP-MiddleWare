package requests

import "lbe/model"

type UserTierUpdateEventRequest struct {
	EventLookup string `json:"event_lookup,omitempty"`
	UserId      string `json:"user_id,omitempty"`
	RetailerID  string `json:"retailer_id,omitempty"`
}

type UserProfileRequest struct {
	User model.RlpUserReq `json:"user"`
}
