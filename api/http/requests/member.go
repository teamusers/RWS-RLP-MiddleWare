package requests

// User represents the profile fields for creating or updating a user.
// All fields are optional when updating—only non‑zero values will be applied.
type MemberUser struct {
	// Email is the user’s email address.
	Email string `json:"email" example:"user@example.com"`

	// GR_ID is the group or partner system identifier for the user.
	GR_ID string `json:"gr_id" example:"GR12345"`

	// RLP_ID is the RLP system identifier for the user.
	RLP_ID string `json:"rlp_id" example:"20250430000001"`

	// RWS_Membership_ID is the RWS membership ID assigned to this user.
	RLP_NO string `json:"rlp_no" example:"70000000001"`
}

// CreateUserRequest wraps the data needed to register a new user.
// It includes both the core User fields and the email for initial signup checks.
type CreateMemberUser struct {
	// User contains the fields to create in the new user record.
	User MemberUser `json:"user"`
}

// UpdateBurnPinRequest is the payload for updating a user’s burn PIN.
type UpdateBurnPin struct {
	// Email of the user whose burn PIN is being updated.
	Email string `json:"email" binding:"required" example:"user@example.com"`

	// BurnPin is the new numeric PIN to set.
	BurnPin int64 `json:"burn_pin" binding:"required" example:"4321"`
}

type VerifyUser struct {
	// Email is the user’s email address for verification and lookup.
	Email string `json:"email" example:"user@example.com"`
}
