package model

type User struct {
	Email           string       `json:"email"`
	Identifier      []Identifier `json:"identifiers"`
	MobileCode      int          `json:"mobile_code"`
	MobileNumber    int          `json:"mobile_number"`
	FirstName       string       `json:"first_name"`
	LastName        string       `json:"last_name"`
	DateOfBirth     Date         `json:"dob"`
	CreatedAt       DateTime     `json:"created_at"`
	Country         string       `json:"country"`
	AvailablePoints int          `json:"available_points"`
	Tier            string       `json:"tier"`
	RegisteredAt    DateTime     `json:"registered_at"`
	Suspended       bool         `json:"suspended"`
	UpdatedAt       DateTime     `json:"updated_at"`
	UserProfile     UserProfile  `json:"user_profile"`
}

type Identifier struct {
	ExternalID     string `json:"external_id"`
	ExternalIDType string `json:"external_id_type"`
}

type UserProfile struct {
	PreviousEmail string `json:"previous_email"`
	ActiveStatus  int    `json:"active_status"`
	Language      string `json:"language"`
	MarketingPreference
}

type MarketingPreference struct {
	Push   bool `json:"market_pref_push"`
	Email  bool `json:"market_pref_email"`
	Mobile bool `json:"market_pref_mobile"`
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
