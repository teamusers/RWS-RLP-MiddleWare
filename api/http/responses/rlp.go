package responses

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Date is a custom type based on time.Time that works with JSON and SQL.
// It serializes to "YYYY-MM-DD" and handles null/empty as JSON null.
//
// Example JSON: "2007-08-05"
type Date time.Time

// UnmarshalJSON parses a JSON string in the "2006-01-02" format,
// gracefully handling empty strings and JSON null.
func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*d = Date(time.Time{})
		return nil
	}
	s := strings.Trim(string(data), `"`)
	if s == "" {
		*d = Date(time.Time{})
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

// MarshalJSON outputs the date in "2006-01-02" format.
// If the date is the zero value, it outputs JSON null.
//
// Example output: "2007-08-05"
func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if t.IsZero() {
		return []byte("null"), nil
	}
	formatted := t.Format("2006-01-02")
	return []byte(fmt.Sprintf(`"%s"`, formatted)), nil
}

// Value implements the driver.Valuer interface.
// Returns nil (SQL NULL) if the Date is zero.
func (d Date) Value() (driver.Value, error) {
	t := time.Time(d)
	if t.IsZero() {
		return nil, nil
	}
	return t, nil
}

// Scan implements the sql.Scanner interface.
// Converts a database value into a Date, handling NULL.
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		*d = Date(time.Time{})
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("failed to scan Date: %v", value)
	}
	*d = Date(t)
	return nil
}

type DateTime time.Time

func (dt *DateTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		*dt = DateTime(time.Time{})
		return nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*dt = DateTime(t)
	return nil
}

// GetUserResponse represents the top-level JSON
type GetUserResponse struct {
	Status string `json:"status"`
	User   User   `json:"user"`
}

// User holds all the user fields
type User struct {
	ID              string        `json:"id"`
	ExternalID      string        `json:"external_id"`
	ProxyIDs        []string      `json:"proxy_ids"`
	OptedIn         bool          `json:"opted_in"`
	Email           string        `json:"email"`
	Identifiers     []Identifier  `json:"identifiers"`
	FirstName       string        `json:"first_name"`
	LastName        string        `json:"last_name"`
	Gender          string        `json:"gender"`
	Dob             Date          `json:"dob"` // format: "2006-01-02"
	AccountStatus   string        `json:"account_status"`
	AuthToken       string        `json:"auth_token"`
	CreatedAt       DateTime      `json:"created_at"` // format: "2006-01-02 15:04:05"
	Address         string        `json:"address"`
	Address2        string        `json:"address2"`
	City            string        `json:"city"`
	State           string        `json:"state"`
	Zip             string        `json:"zip"`
	Country         string        `json:"country"`
	AvailablePoints int           `json:"available_points"`
	Tier            string        `json:"tier"`
	ReferrerCode    string        `json:"referrer_code"`
	RegisteredAt    DateTime      `json:"registered_at"` // same format as CreatedAt
	Suspended       bool          `json:"suspended"`
	UpdatedAt       DateTime      `json:"updated_at"` // same format as CreatedAt
	PhoneNumbers    []PhoneNumber `json:"phone_numbers"`
}

// Identifier represents an external ID mapping
type Identifier struct {
	ExternalID     string `json:"external_id"`
	ExternalIDType string `json:"external_id_type"`
}

// PhoneNumber holds a phone record
type PhoneNumber struct {
	PhoneNumber       string   `json:"phone_number"`
	PhoneType         string   `json:"phone_type"`
	PreferenceFlags   []string `json:"preference_flags"`
	VerifiedOwnership bool     `json:"verified_ownership"`
}
