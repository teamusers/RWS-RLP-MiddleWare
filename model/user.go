package model

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

// User represents a user account.
// swagger:model User
type User struct {
	// ID is the auto-incrementing primary key.
	// example: 1
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement"`

	// AccountStatus indicates the user's account status.
	// example: good
	AccountStatus string `json:"account_status" example:"good"`

	// CreatedAt is the timestamp when the user account was created.
	// example: 2016-10-21T18:12:22Z
	CreatedAt time.Time `json:"created_at" example:"2016-10-21T18:12:22Z"`

	// DOB is the user's date of birth.
	// example: 1980-01-01
	DOB Date `json:"dob" example:"1980-01-01"`

	// Email is the user's email address.
	// example: john.smith@fake.email.addr
	Email string `json:"email" example:"john.smith@fake.email.addr"`

	// ExternalID is the user's external identifier.
	// example: 654321
	ExternalID string `json:"external_id" example:"654321"`

	// OptedIn indicates whether the user has opted in.
	// example: true
	OptedIn bool `json:"opted_in" example:"true"`

	// Gender of the user (m/f).
	// example: m
	Gender string `json:"gender" example:"m"`

	// RegisteredAt is the timestamp when the user completed registration.
	// example: 2016-10-21T18:12:22Z
	RegisteredAt time.Time `json:"registered_at" example:"2016-10-21T18:12:22Z"`

	// Suspended indicates whether the account is suspended.
	// example: false
	Suspended bool `json:"suspended" example:"false"`

	// UpdatedAt is the timestamp of the last update.
	// example: 2016-10-21T18:12:22Z
	UpdatedAt time.Time `json:"updated_at" example:"2016-10-21T18:12:22Z"`

	// ReferrerCode is the code used at sign‑up.
	// example: JOHN-70A756
	ReferrerCode string `json:"referrer_code" example:"JOHN-70A756"`

	// PhoneNumbers associated with this user.
	PhoneNumbers []PhoneNumber `json:"phone_numbers" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// Identifiers from external systems.
	Identifiers []Identifier `json:"identifiers" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// TableName explicitly sets the SQL table name for User.
func (User) TableName() string {
	return "users"
}

// PhoneNumber holds one of the user’s phone numbers.
// swagger:model PhoneNumber
type PhoneNumber struct {
	ID                int64  `json:"-" gorm:"primaryKey;autoIncrement"`
	UserID            int64  `json:"-" gorm:"index"` // ← foreign key back to users.id
	PhoneNumber       string `json:"phone_number" example:"1234123123"`
	PhoneType         string `json:"phone_type" example:"home"`
	VerifiedOwnership bool   `json:"verified_ownership" example:"false"`
}

// TableName sets the SQL table name for UserPhoneNumber.
func (PhoneNumber) TableName() string {
	return "users_phone_numbers"
}

// Identifier holds an external ID and its type.
// swagger:model Identifier
type Identifier struct {
	ID             int64  `json:"-" gorm:"primaryKey;autoIncrement"`
	UserID         int64  `json:"-" gorm:"index"`
	ExternalID     string `json:"external_id" example:"1234abcd"`
	ExternalIDType string `json:"external_id_type" example:"facebook"`
}

// TableName sets the SQL table name for UserPhoneNumber.
func (Identifier) TableName() string {
	return "users_identifiers"
}
