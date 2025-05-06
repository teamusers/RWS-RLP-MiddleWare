package responses

import (
	"lbe/model"
)

// GetUserResponse represents the top-level JSON
type GetUserResponse struct {
	Status string            `json:"status"`
	User   model.RlpUserResp `json:"user"`
}

// // User holds all the user fields
// // TO DO: Refactor to RlpUser
// type User struct {
// 	ID              string             `json:"id"`
// 	ExternalID      string             `json:"external_id"`
// 	ProxyIDs        []string           `json:"proxy_ids"`
// 	OptedIn         bool               `json:"opted_in"`
// 	Email           string             `json:"email"`
// 	Identifiers     []model.Identifier `json:"identifiers"` //TO DO: Standardize
// 	FirstName       string             `json:"first_name"`
// 	LastName        string             `json:"last_name"`
// 	Gender          string             `json:"gender"`
// 	Dob             model.Date         `json:"dob"` // format: "2006-01-02"
// 	AccountStatus   string             `json:"account_status"`
// 	AuthToken       string             `json:"auth_token"`
// 	CreatedAt       model.DateTime     `json:"created_at"` // format: "2006-01-02 15:04:05"
// 	Address         string             `json:"address"`
// 	Address2        string             `json:"address2"`
// 	City            string             `json:"city"`
// 	State           string             `json:"state"`
// 	Zip             string             `json:"zip"`
// 	Country         string             `json:"country"`
// 	AvailablePoints int                `json:"available_points"`
// 	Tier            string             `json:"tier"`
// 	ReferrerCode    string             `json:"referrer_code"`
// 	RegisteredAt    model.DateTime     `json:"registered_at"` // same format as CreatedAt
// 	Suspended       bool               `json:"suspended"`
// 	UpdatedAt       model.DateTime     `json:"updated_at"` // same format as CreatedAt
// 	PhoneNumbers    []PhoneNumber      `json:"phone_numbers"`
// }

// // // Identifier represents an external ID mapping
// // type Identifier struct {
// // 	ExternalID     string `json:"external_id"`
// // 	ExternalIDType string `json:"external_id_type"`
// // }

// // PhoneNumber holds a phone record
// type PhoneNumber struct {
// 	PhoneNumber       string   `json:"phone_number"`
// 	PhoneType         string   `json:"phone_type"`
// 	PreferenceFlags   []string `json:"preference_flags"`
// 	VerifiedOwnership bool     `json:"verified_ownership"`
// }
