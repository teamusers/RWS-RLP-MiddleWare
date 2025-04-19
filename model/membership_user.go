package model

import "time"

// MembershipUser represents a full user profile in the membership system.
//
// Example:
// {
//   "id": 123,
//   "external_id": "abc123",
//   "external_id_type": "EMAIL",
//   "email": "user@example.com",
//   "burn_pin": 4321,
//   "session_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
//   "session_expiry": 1744176000,
//   "gr_id": "GR12345",
//   "rlp_id": "RLP67890",
//   "rws_membership_id": "RWS54321",
//   "rws_membership_number": 1000,
//   "created_at": "2025-04-19T10:00:00Z",
//   "updated_at": "2025-04-19T11:00:00Z"
// }
type MembershipUser struct {
	// ID is the database primary key.
	// example: 123
	ID uint64 `json:"id" example:"123"`

	// ExternalID is the client‑side unique identifier for the user.
	// example: "abc123"
	ExternalID string `json:"external_id" example:"abc123"`

	// ExternalTYPE describes the kind/source of the external ID.
	// example: "EMAIL"
	ExternalTYPE string `json:"external_id_type" example:"EMAIL"`

	// Email of the user.
	// example: "user@example.com"
	Email string `json:"email" example:"user@example.com"`

	// BurnPin is the numeric PIN used for burn operations.
	// example: 4321
	BurnPin uint64 `json:"burn_pin" example:"4321"`

	// SessionToken is the JWT used for subsequent API calls.
	// example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	SessionToken string `json:"session_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`

	// SessionExpiry is the Unix timestamp when SessionToken expires.
	// example: 1744176000
	SessionExpiry int64 `json:"session_expiry" example:"1744176000"`

	// GR_ID is the group‑related ID in the GR system.
	// example: "GR12345"
	GR_ID string `json:"gr_id" example:"GR12345"`

	// RLP_ID is the RLP system identifier for this user.
	// example: "RLP67890"
	RLP_ID string `json:"rlp_id" example:"RLP67890"`

	// RWS_Membership_ID is the RWS membership ID assigned by the loyalty platform.
	// example: "RWS54321"
	RWS_Membership_ID string `json:"rws_membership_id" example:"RWS54321"`

	// RWS_Membership_Number is the numeric membership number in RWS.
	// example: 1000
	RWS_Membership_Number uint64 `json:"rws_membership_number" example:"1000"`

	// CreatedAt is the timestamp when the record was first created.
	// example: "2025-04-19T10:00:00Z"
	CreatedAt time.Time `json:"created_at" example:"2025-04-19T10:00:00Z"`

	// UpdatedAt is the timestamp when the record was last updated.
	// example: "2025-04-19T11:00:00Z"
	UpdatedAt time.Time `json:"updated_at" example:"2025-04-19T11:00:00Z"`
}
