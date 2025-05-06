package model

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
	// example: premium
	Class string `json:"class,omitempty" example:"premium"`

	// User’s first name
	// example: Alice
	FirstName string `json:"first_name,omitempty" example:"Alice"`

	// User’s last name
	// example: Smith
	LastName string `json:"last_name,omitempty" example:"Smith"`

	// User’s email address
	// example: alice.smith@example.com
	Email string `json:"email,omitempty" example:"alice.smith@example.com"`

	// Date of birth (YYYY-MM-DD)
	// example: 1985-07-20
	DateOfBirth Date `json:"dob,omitempty" example:"1985-07-20"`

	// Country dialing code
	// example: 65
	MobileCode int `json:"mobile_code,omitempty" example:"65"`

	// Local mobile number
	// example: 91234567
	MobileNumber int `json:"mobile_number,omitempty" example:"91234567"`
}

func (gr *GrProfile) MapGrProfileToLbeUser() User {
	return User{
		Email:        gr.Email,
		FirstName:    gr.FirstName,
		LastName:     gr.LastName,
		DateOfBirth:  gr.DateOfBirth,
		MobileCode:   gr.MobileCode,
		MobileNumber: gr.MobileNumber,
	}
}
