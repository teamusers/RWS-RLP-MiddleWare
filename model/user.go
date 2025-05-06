package model

// User represents a customer in the system.
// swagger:model User
type User struct {
	// Email address of the user
	// example: john.doe@example.com
	Email string `json:"email" example:"john.doe@example.com"`

	// List of external identifiers for the user
	// example: [{"external_id":"ABC123","external_id_type":"loyalty"}]
	Identifier []Identifier `json:"identifiers"`

	// Country dialing code
	// example: 1
	MobileCode int `json:"mobile_code" example:"1"`

	// Mobile phone number
	// example: 98765432
	MobileNumber int `json:"mobile_number" example:"98765432"`

	// User's first name
	// example: John
	FirstName string `json:"first_name" example:"John"`

	// User's last name
	// example: Doe
	LastName string `json:"last_name" example:"Doe"`

	// Date of birth in YYYY-MM-DD
	// example: 1990-05-15
	DateOfBirth Date `json:"dob" example:"1990-05-15"`

	// Timestamp when the record was created (RFC3339)
	// example: 2025-04-01T08:30:00Z
	CreatedAt DateTime `json:"created_at" example:"2006-01-02 15:04:05"`

	// ISO 3166-1 alpha-2 country code
	// example: SG
	Country string `json:"country" example:"SG"`

	// Loyalty points available
	// example: 1200
	AvailablePoints int `json:"available_points" example:"1200"`

	// Loyalty tier name
	// example: gold
	Tier string `json:"tier" example:"gold"`

	// Timestamp when the user registered
	// example: 2025-04-01T08:30:00Z
	RegisteredAt DateTime `json:"registered_at" example:"2006-01-02 15:04:05"`

	// Whether the account is suspended
	// example: false
	Suspended bool `json:"suspended" example:"false"`

	// Timestamp of last update
	// example: 2025-05-05T14:00:00Z
	UpdatedAt DateTime `json:"updated_at" example:"2006-01-02 15:04:05"`

	// Additional profile details
	UserProfile UserProfile `json:"user_profile"`
}

// Identifier holds an external ID and its type.
// swagger:model Identifier
type Identifier struct {
	// The external identifier value
	// example: ABC123
	ExternalID string `json:"external_id" example:"ABC123"`

	// Type of the external identifier
	// example: loyalty
	ExternalIDType string `json:"external_id_type" example:"loyalty"`
}

// UserProfile contains supplementary user attributes.
// swagger:model UserProfile
type UserProfile struct {
	// Previously used email
	// example: john.old@example.com
	PreviousEmail string `json:"previous_email" example:"john.old@example.com"`

	// Active status code (e.g., 1=active, 0=inactive)
	// example: 1
	ActiveStatus int `json:"active_status" example:"1"`

	// Preferred language (ISO 639-1)
	// example: en
	Language string `json:"language" example:"en"`

	// Embedded marketing preferences
	MarketingPreference
}

// MarketingPreference defines the user's marketing opt-in channels.
// swagger:model MarketingPreference
type MarketingPreference struct {
	// Whether the user opts in to push notifications
	// example: true
	Push bool `json:"market_pref_push" example:"true"`

	// Whether the user opts in to email marketing
	// example: false
	Email bool `json:"market_pref_email" example:"false"`

	// Whether the user opts in to SMS/mobile marketing
	// example: true
	Mobile bool `json:"market_pref_mobile" example:"true"`
}

func (u *User) MapLbeToRlpUser(rlpId string) RlpUserReq {
	return RlpUserReq{
		ExternalID:  rlpId,
		Email:       u.Email,
		Dob:         u.DateOfBirth,
		Country:     u.Country,
		UserProfile: u.UserProfile,
	}
	// TO DO / TBC : add identifiers
}
