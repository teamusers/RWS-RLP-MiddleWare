package requests

// User represents the profile fields for creating or updating a user.
// All fields are optional when updating—only non‑zero values will be applied.
type User struct {
	// ExternalID is the client system’s unique identifier for this user.
	ExternalID string `json:"external_id" example:"abc123"`

	// ExternalTYPE describes the type or source of the external ID.
	ExternalTYPE string `json:"external_id_type" example:"EMAIL"`

	// Email is the user’s email address.
	Email string `json:"email" example:"user@example.com"`

	// BurnPin is the numeric PIN used for burn operations.
	BurnPin uint64 `json:"burn_pin" example:"1234"`

	// SessionToken is the login session token issued to the user.
	SessionToken string `json:"session_token" example:"eyJhbGciOiJIUzI1..."`

	// SessionExpiry is the Unix timestamp (seconds since epoch) when the session token expires.
	SessionExpiry int64 `json:"session_expiry" example:"1712345678"`

	// GR_ID is the group or partner system identifier for the user.
	GR_ID string `json:"gr_id" example:"GR12345"`

	// RLP_ID is the RLP system identifier for the user.
	RLP_ID string `json:"rlp_id" example:"RLP67890"`

	// RWS_Membership_ID is the RWS membership ID assigned to this user.
	RWS_Membership_ID string `json:"rws_membership_id" example:"RWS54321"`

	// RWS_Membership_Number is the numeric membership number in the RWS system.
	RWS_Membership_Number uint64 `json:"rws_membership_number" example:"987654"`
}

// UpdateBurnPinRequest is the payload for updating a user’s burn PIN.
type UpdateBurnPin struct {
	// Email of the user whose burn PIN is being updated.
	Email string `json:"email" binding:"required" example:"user@example.com"`

	// BurnPin is the new numeric PIN to set.
	BurnPin int64 `json:"burn_pin" binding:"required" example:"4321"`
}

// CreateUserRequest wraps the data needed to register a new user.
// It includes both the core User fields and the email for initial signup checks.
type CreateUser struct {
	// User contains the fields to create in the new user record.
	User User `json:"user"`
}
type VerifyUser struct {
	// Email is the user’s email address for verification and lookup.
	Email string `json:"email" example:"user@example.com"`
}
