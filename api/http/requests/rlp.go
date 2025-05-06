package requests

// RlpUser represents the payload to register or update a customer.
// swagger:model RlpUser
// type RlpUser struct {
// 	// Identifier for customer in external system.
// 	// Required if Email is not specified.
// 	ExternalID string `json:"external_id" example:"1284111"`

// 	// Indicates opt-in to loyalty program (defaults to true if not set).
// 	OptedIn *bool `json:"opted_in,omitempty" example:"true"`

// 	// Type associated with external identifier.
// 	ExternalIDType string `json:"external_id_type,omitempty" example:"facebook"`

// 	// Customer’s email (encrypted at rest). Required if ExternalID is not specified.
// 	Email string `json:"email,omitempty" example:"john.smith@test.com"`

// 	// Preferred locale, e.g. en-US
// 	Locale string `json:"locale,omitempty" example:"en-US"`

// 	// Customer’s IP at registration.
// 	IP string `json:"ip,omitempty" example:"203.0.113.42"`

// 	// Date of birth, YYYY-MM-DD
// 	Dob model.Date `json:"dob,omitempty" example:"1980-01-01"`

// 	// Address line 1
// 	Address string `json:"address,omitempty" example:"7 Tremont Street"`

// 	// Address line 2
// 	Address2 string `json:"address2,omitempty" example:"8 Tremont Street"`

// 	// City of residence
// 	City string `json:"city,omitempty" example:"Boston"`

// 	// State/province/region
// 	State string `json:"state,omitempty" example:"MA"`

// 	// Zip or postal code
// 	Zip string `json:"zip,omitempty" example:"02021"`

// 	// 3-letter ISO-3166-1 country code
// 	Country string `json:"country,omitempty" example:"USA"`

// 	// Gender of customer: m or f
// 	Gender string `json:"gender,omitempty" example:"m"`

// 	// First name
// 	FirstName string `json:"first_name,omitempty" example:"John"`

// 	// Last name (surname)
// 	LastName string `json:"last_name,omitempty" example:"Smith"`

// 	// Phone numbers associated with customer
// 	PhoneNumbers []PhoneNumber `json:"phone_numbers,omitempty"`

// 	// Referral info when signing up with a referrer code
// 	Referral *Referral `json:"referral,omitempty"`
// }

// // PhoneNumber holds a single phone entry.
// // swagger:model PhoneNumber
// type PhoneNumber struct {
// 	// Phone number (digits only)
// 	PhoneNumber string `json:"phone_number" example:"1234123123"`

// 	// Type of phone number: mobile, office, home, fax, other
// 	PhoneType string `json:"phone_type,omitempty" example:"home"`

// 	// Flags like "primary"
// 	PreferenceFlags []string `json:"preference_flags,omitempty" example:"[\"primary\"]"`
// }

// // Referral holds referral info.
// // swagger:model Referral
// type Referral struct {
// 	// Code of the referring customer
// 	ReferralCode string `json:"referral_code,omitempty" example:"JOHN-70A756"`
// }
