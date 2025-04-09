package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Date is a custom type based on time.Time that works with JSON and SQL.
type Date time.Time

// UnmarshalJSON parses a JSON string in the "2006-01-02" format,
// while gracefully handling empty strings and JSON null.
func (d *Date) UnmarshalJSON(data []byte) error {
	// Check if the JSON value is null.
	if string(data) == "null" {
		*d = Date(time.Time{})
		return nil
	}

	// Remove the surrounding quotes.
	s := strings.Trim(string(data), "\"")
	// Return zero value if the string is empty.
	if s == "" {
		*d = Date(time.Time{})
		return nil
	}

	// Parse the non-empty date string.
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

// MarshalJSON outputs the date in "2006-01-02" format.
// If the date is the zero value, it outputs JSON null.
func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if t.IsZero() {
		return []byte("null"), nil
	}
	formatted := t.Format("2006-01-02")
	return []byte(fmt.Sprintf("\"%s\"", formatted)), nil
}

// Value implements the driver.Valuer interface.
// It returns nil (SQL NULL) if the Date is the zero value.
func (d Date) Value() (driver.Value, error) {
	t := time.Time(d)
	if t.IsZero() {
		return nil, nil
	}
	return t, nil
}

// Scan implements the sql.Scanner interface.
// It converts a database value into a Date, handling NULL values.
func (d *Date) Scan(value interface{}) error {
	// If the DB column is NULL, assign zero value.
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

// User represents a user model that maps to a MySQL table.
type User struct {
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ExternalID   string    `gorm:"column:external_id;size:50" json:"external_id"`
	OptedIn      bool      `gorm:"column:opted_in" json:"opted_in"`
	ExternalTYPE string    `gorm:"column:external_id_type;size:50" json:"external_id2"`
	Email        string    `gorm:"column:email;size:255" json:"email"`
	DOB          Date      `gorm:"column:dob" json:"dob"`
	Country      string    `gorm:"column:country;size:3" json:"country"`
	FirstName    string    `gorm:"column:first_name;size:255" json:"first_name"`
	LastName     string    `gorm:"column:last_name;size:255" json:"last_name"`
	BurnPin      string    `gorm:"column:burn_pin;size:4" json:"burn_pin"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`

	// PhoneNumbers represents the one-to-many relationship to UserPhoneNumber.
	PhoneNumbers []UserPhoneNumber `gorm:"foreignKey:UserID;references:ID" json:"phone_numbers"`
}

// TableName sets the table name for the User model.
func (User) TableName() string {
	return "users"
}

type UserPhoneNumber struct {
	ID              uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID          uint64    `gorm:"column:user_id" json:"user_id"` // foreign key to User.ID
	PhoneNumber     string    `gorm:"column:phone_number;size:20" json:"phone_number"`
	PhoneType       string    `gorm:"column:phone_type;size:20" json:"phone_type"`
	PreferenceFlags string    `gorm:"column:preference_flags;size:50" json:"preference_flags"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// If your actual table name is `user_phone_numbers` (plural), you can
// explicitly specify it below. Otherwise, GORM might pluralize the struct
// name by default.
func (UserPhoneNumber) TableName() string {
	return "users_phone_numbers"
}
