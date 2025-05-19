package requests

type UserTierUpdateEventRequest struct {
	EventLookup string `json:"event_lookup,omitempty"`
	UserId      string `json:"user_id,omitempty"`
	UserTier    string `json:"user_tier,omitempty"`
}
