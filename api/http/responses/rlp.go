package responses

import (
	"lbe/model"
	"time"
)

// GetUserResponse represents the top-level JSON
type GetUserResponse struct {
	Status string            `json:"status"`
	User   model.RlpUserResp `json:"user"`
}

type UserTierUpdateEventResponse struct {
	JsonBody UserTierUpdateJsonBody `json:"application/json"`
}

type UserTierUpdateJsonBody struct {
	UserEvent                  UserTierUpdateUserEvent `json:"user_event"`
	Culture                    string                  `json:"culture"`
	EventSavedSuccessfully     bool                    `json:"event_saved_successfully"`
	RulesProcessedSuccessfully bool                    `json:"rules_processed_successfully"`
	Outcomes                   []UserTierUpdateOutcome `json:"outcomes"`
}

type UserTierUpdateUserEvent struct {
	ID               string    `json:"id"`
	TimeOfOccurrence time.Time `json:"time_of_occurrence"`
	RetailerID       string    `json:"retailer_id"`
	EventLookup      string    `json:"event_lookup"`
	UserID           string    `json:"user_id"`
	UserTier         string    `json:"user_tier"`
	IsSessionM       bool      `json:"is_session_m"`
}

type UserTierUpdateOutcome struct {
	IsOutcomeApplied bool   `json:"is_outcome_applied"`
	OutcomeMessage   string `json:"outcome_message"`
	UserTier         string `json:"user_tier"`
}
