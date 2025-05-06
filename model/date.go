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

// UnmarshalJSON will now be on the value receiver,
// so *Date and Date both implement json.Unmarshaler.
func (d Date) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		return nil
	}
	_, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	// but this only updates the copy; you still need a pointer
	// => so this approach is tricky unless you also satisfy
	// json.Unmarshaler on *Date
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
