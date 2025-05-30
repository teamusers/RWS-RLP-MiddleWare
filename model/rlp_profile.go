package model

// RlpUser represents the payload to register or update a customer.
// swagger:model RlpUserReq
type RlpUserReq struct {
	// Identifier for customer in external system.
	// Required if Email is not specified.
	ExternalID string `json:"external_id,omitempty" example:"1284111"`

	// Indicates opt-in to loyalty program (defaults to true if not set).
	OptedIn bool `json:"opted_in,omitempty" example:"true"`

	// Type associated with external identifier.
	ExternalIDType string `json:"external_id_type,omitempty" example:"facebook"`

	// Array of external id identifiers tied to the user.
	Identifier []Identifier `json:"identifiers,omitempty"`

	// Customer’s email (encrypted at rest). Required if ExternalID is not specified.
	Email string `json:"email,omitempty" example:"john.smith@test.com"`

	// Preferred locale, e.g. en-US
	Locale string `json:"locale,omitempty" example:"en-US"`

	// Customer’s IP at registration.
	IP string `json:"ip,omitempty" example:"203.0.113.42"`

	// Date of birth, YYYY-MM-DD
	Dob *Date `json:"dob,omitempty" example:"1980-01-01"`

	// Address line 1
	Address string `json:"address,omitempty" example:"7 Tremont Street"`

	// Address line 2
	Address2 string `json:"address2,omitempty" example:"8 Tremont Street"`

	// City of residence
	City string `json:"city,omitempty" example:"Boston"`

	// State/province/region
	State string `json:"state,omitempty" example:"MA"`

	// Zip or postal code
	Zip string `json:"zip,omitempty" example:"02021"`

	// 3-letter ISO-3166-1 country code
	Country string `json:"country,omitempty" example:"USA"`

	// Gender of customer: m or f
	Gender string `json:"gender,omitempty" example:"m"`

	// First name
	FirstName string `json:"first_name,omitempty" example:"John"`

	// Last name (surname)
	LastName string `json:"last_name,omitempty" example:"Smith"`

	// Phone numbers associated with customer
	PhoneNumbers []PhoneNumber `json:"phone_numbers,omitempty"`

	// Referral info when signing up with a referrer code
	Referral *Referral `json:"referral,omitempty"`

	// User profile for custom attributes
	UserProfile UserProfile `json:"user_profile,omitempty"`
}

// Referral holds referral info.
// swagger:model Referral
type Referral struct {
	// Code of the referring customer
	ReferralCode string `json:"referral_code,omitempty" example:"JOHN-70A756"`
}

// RlpUserResp holds all the user fields
// swagger:model RlpUserResp
type RlpUserResp struct {
	ID              string        `json:"id"`
	ExternalID      string        `json:"external_id"`
	ProxyIDs        []string      `json:"proxy_ids"`
	OptedIn         bool          `json:"opted_in"`
	Email           string        `json:"email"`
	Identifiers     []Identifier  `json:"identifiers"`
	FirstName       string        `json:"first_name"`
	LastName        string        `json:"last_name"`
	Gender          string        `json:"gender"`
	Dob             *Date         `json:"dob"` // format: "2006-01-02"
	AccountStatus   string        `json:"account_status"`
	AuthToken       string        `json:"auth_token"`
	CreatedAt       *DateTime     `json:"created_at"` // format: "2006-01-02 15:04:05"
	Address         string        `json:"address"`
	Address2        string        `json:"address2"`
	City            string        `json:"city"`
	State           string        `json:"state"`
	Zip             string        `json:"zip"`
	Country         string        `json:"country"`
	AvailablePoints float64       `json:"available_points"`
	Tier            string        `json:"tier"`
	ReferrerCode    string        `json:"referrer_code"`
	RegisteredAt    *DateTime     `json:"registered_at"` // same format as CreatedAt
	Suspended       bool          `json:"suspended"`
	UpdatedAt       *DateTime     `json:"updated_at"` // same format as CreatedAt
	PhoneNumbers    []PhoneNumber `json:"phone_numbers"`
	UserProfile     UserProfile   `json:"user_profile"` //TBC
}

func (rlpUser *RlpUserResp) MapRlpToLbeUser() User {
	return User{
		Email:           rlpUser.Email,
		Identifier:      rlpUser.Identifiers,
		PhoneNumbers:    rlpUser.PhoneNumbers,
		FirstName:       rlpUser.FirstName,
		LastName:        rlpUser.LastName,
		DateOfBirth:     rlpUser.Dob,
		CreatedAt:       rlpUser.CreatedAt,
		Country:         rlpUser.Country,
		AvailablePoints: rlpUser.AvailablePoints,
		Tier:            rlpUser.Tier,
		RegisteredAt:    rlpUser.RegisteredAt,
		Suspended:       rlpUser.Suspended,
		UpdatedAt:       rlpUser.UpdatedAt,
		UserProfile:     rlpUser.UserProfile,
	}
}

func (u *RlpUserReq) PopulateRegistrationDefaults(rlpId string) {
	u.OptedIn = true
	u.ExternalID = rlpId
	u.ExternalIDType = "RLP_ID"
	u.UserProfile.LanguagePreference = "EN"
	u.UserProfile.PreviousEmail = u.Email
	u.UserProfile.MarketingPreference.Email = BoolPtr(true)
	u.UserProfile.MarketingPreference.Push = BoolPtr(true)
	u.UserProfile.MarketingPreference.Mobile = BoolPtr(true)
	u.UserProfile.ActiveStatus = "1"
}

func BoolPtr(b bool) *bool {
	return &b
}
