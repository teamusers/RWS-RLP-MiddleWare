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

// User represents a customer record in the users table.
//
// Example JSON:
//
//	{
//	  "id": 42,
//	  "external_id": "abc123",
//	  "opted_in": true,
//	  "external_id_type": "EMAIL",
//	  "email": "user@example.com",
//	  "dob": "2007-08-05",
//	  "country": "SGP",
//	  "first_name": "Brendan",
//	  "last_name": "Test",
//	  "burn_pin": 1234,
//	  "created_at": "2025-04-19T10:00:00Z",
//	  "updated_at": "2025-04-19T11:00:00Z",
//	  "phone_numbers": [ … ]
//	}
type User struct {
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id" example:"42"`
	ExternalID   string    `gorm:"column:external_id;size:50" json:"external_id" example:"abc123"`
	OptedIn      bool      `gorm:"column:opted_in" json:"opted_in" example:"true"`
	ExternalTYPE string    `gorm:"column:external_id_type;size:50" json:"external_id_type" example:"EMAIL"`
	Email        string    `gorm:"column:email;size:255" json:"email" example:"user@example.com"`
	DOB          Date      `gorm:"column:dob" json:"dob" example:"2007-08-05"`
	Country      string    `gorm:"column:country;size:3" json:"country" example:"SGP"`
	FirstName    string    `gorm:"column:first_name;size:255" json:"first_name" example:"Brendan"`
	LastName     string    `gorm:"column:last_name;size:255" json:"last_name" example:"Test"`
	BurnPin      uint64    `gorm:"column:burn_pin;size:4" json:"burn_pin" example:"1234"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at" example:"2025-04-19T10:00:00Z"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at" example:"2025-04-19T11:00:00Z"`

	// PhoneNumbers holds zero or more phone numbers associated with this user.
	PhoneNumbers []UserPhoneNumber `gorm:"foreignKey:UserID;references:ID" json:"phone_numbers"`
}

// TableName explicitly sets the SQL table name for User.
func (User) TableName() string {
	return "users"
}

// UserPhoneNumber represents a phone number linked to a user.
//
// Example JSON:
//
//	{
//	  "id": 101,
//	  "user_id": 42,
//	  "phone_number": "+6598765432",
//	  "phone_type": "mobile",
//	  "preference_flags": "primary",
//	  "created_at": "2025-04-19T10:05:00Z",
//	  "updated_at": "2025-04-19T10:05:00Z"
//	}
type UserPhoneNumber struct {
	ID              uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id" example:"101"`
	UserID          uint64    `gorm:"column:user_id" json:"user_id" example:"42"`
	PhoneNumber     string    `gorm:"column:phone_number;size:20" json:"phone_number" example:"+6598765432"`
	PhoneType       string    `gorm:"column:phone_type;size:20" json:"phone_type" example:"mobile"`
	PreferenceFlags string    `gorm:"column:preference_flags;size:50" json:"preference_flags" example:"primary"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at" example:"2025-04-19T10:05:00Z"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at" example:"2025-04-19T10:05:00Z"`
}

// TableName sets the SQL table name for UserPhoneNumber.
func (UserPhoneNumber) TableName() string {
	return "users_phone_numbers"
}
