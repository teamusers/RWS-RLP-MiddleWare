package model

// User represents a customer in the system.
// swagger:model User
type User struct {
	// Email address of the user
	// example: john.doe@example.com
	Email string `json:"email,omitempty" example:"john.doe@example.com"`

	// List of external identifiers for the user
	// example: [{"external_id":"25052300047","external_id_type":"rlp_id"}]
	Identifier []Identifier `json:"identifiers,omitempty"`

	// Mobile phone number array
	// example: [{"phone_number":"87654321"}]
	PhoneNumbers []PhoneNumber `json:"phone_numbers,omitempty"`

	// User's first name
	// example: John
	FirstName string `json:"first_name,omitempty" example:"John"`

	// User's last name
	// example: Doe
	LastName string `json:"last_name,omitempty" example:"Doe"`

	// Date of birth in YYYY-MM-DD
	// example: 1990-05-15
	DateOfBirth *Date `json:"dob,omitempty" example:"1990-05-15"`

	// Timestamp when the record was created (RFC3339)
	// example: 2025-04-01T08:30:00Z
	CreatedAt *DateTime `json:"created_at,omitempty" example:"2006-01-02 15:04:05"`

	// ISO 3166-1 alpha-2 country code
	// example: SG
	Country string `json:"country,omitempty" example:"SG"`

	// Loyalty points available
	// example: 1200
	AvailablePoints int `json:"available_points,omitempty" example:"1200"`

	// Loyalty tier name
	// example: gold
	Tier string `json:"tier,omitempty" example:"gold"`

	// Timestamp when the user registered
	// example: 2025-04-01T08:30:00Z
	RegisteredAt *DateTime `json:"registered_at,omitempty" example:"2006-01-02 15:04:05"`

	// Whether the account is suspended
	// example: false
	Suspended bool `json:"suspended,omitempty" example:"false"`

	// Timestamp of last update
	// example: 2025-05-05T14:00:00Z
	UpdatedAt *DateTime `json:"updated_at,omitempty" example:"2006-01-02 15:04:05"`

	// Additional profile details
	UserProfile UserProfile `json:"user_profile,omitempty"`

	// GR Profile-unique information. Only used within LBE.
	GrProfile *GrProfile `json:"gr_profile,omitempty"`
}

// Identifier holds an external ID and its type.
// swagger:model Identifier
type Identifier struct {
	// The external identifier value
	// example: 25052300047
	ExternalID string `json:"external_id" example:"25052300047"`

	// Type of the external identifier
	// example: rlp_id
	ExternalIDType string `json:"external_id_type" example:"rlp_id"`
}

// PhoneNumber holds a phone record
// swagger:model PhoneNumber
type PhoneNumber struct {
	PhoneNumber       string   `json:"phone_number"`
	PhoneType         string   `json:"phone_type"`
	PreferenceFlags   []string `json:"preference_flags"`
	VerifiedOwnership bool     `json:"verified_ownership,omitempty"`
}

// UserProfile contains supplementary user attributes.
// swagger:model UserProfile
type UserProfile struct {
	// Country code for mobile number
	// example: +65
	CountryCode string `json:"country_code,omitempty" example:"+65"`

	// Country name
	// example: Singapore
	CountryName string `json:"country_name,omitempty" example:"Singapore"`

	// Previously used email
	// example: john.old@example.com
	PreviousEmail string `json:"previous_email,omitempty" example:"john.old@example.com"`

	// Active status code (e.g., 1=active, 0=inactive)
	// example: 1
	ActiveStatus *int `json:"active_status,omitempty" example:"1"`

	// Preferred language (ISO 639-1)
	// example: en
	LanguagePreference string `json:"language_preference,omitempty" example:"en"`

	// Embedded marketing preferences
	MarketingPreference

	// Secret Key for burn transaction
	// example: 1111
	BurnPin string `json:"burn_pin,omitempty" example:"1111"`

	// Employee Number for RWS employees only, otherwise empty
	// example: 1111
	EmployeeNumber string `json:"employee_number,omitempty" example:"1111"`
}

// MarketingPreference defines the user's marketing opt-in channels.
// swagger:model MarketingPreference
type MarketingPreference struct {
	// Whether the user opts in to push notifications
	// example: true
	Push *bool `json:"market_pref_push,omitempty" example:"true"`

	// Whether the user opts in to email marketing
	// example: false
	Email *bool `json:"market_pref_email,omitempty" example:"false"`

	// Whether the user opts in to SMS/mobile marketing
	// example: true
	Mobile *bool `json:"market_pref_mobile,omitempty" example:"true"`
}

// GrProfile represents a user’s profile in the GR system.
// swagger:model GrProfile
type GrProfile struct {
	// Unique identifier for the profile
	// example: 123e4567-e89b-12d3-a456-426614174000
	Id string `json:"id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`

	// Four-digit PIN for quick auth
	// example: 1234
	Pin string `json:"pin,omitempty" example:"1234"`

	// User’s membership class
	// example: 1
	Class string `json:"class,omitempty" example:"1"`
}

// Mapper function to convert LBE User format to RLP User format
func (u *User) MapLbeToRlpUser() RlpUserReq {
	return RlpUserReq{
		Identifier:   u.Identifier,
		Email:        u.Email,
		Dob:          u.DateOfBirth,
		Country:      u.Country,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		PhoneNumbers: u.PhoneNumbers,
		UserProfile:  u.UserProfile,
	}
}

// Populate function used to append id links to identifier list for initial registration
func (u *User) PopulateIdentifiers(rlpId, rlpNo string) {
	u.Identifier = append(u.Identifier,
		Identifier{
			ExternalID:     rlpId,
			ExternalIDType: "rlp_id",
		},
		Identifier{
			ExternalID:     rlpNo,
			ExternalIDType: "rlp_no",
		})

	// add gr_id only if in use
	if u.GrProfile != nil && u.GrProfile.Id != "" {
		u.Identifier = append(u.Identifier, Identifier{
			ExternalID:     u.GrProfile.Id,
			ExternalIDType: "gr_id",
		})
	}
}
