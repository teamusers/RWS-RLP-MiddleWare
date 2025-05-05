package model

type GrProfile struct {
	Id           string `json:"id,omitempty"`
	Pin          string `json:"pin,omitempty"`
	Class        string `json:"class,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Email        string `json:"email,omitempty"`
	DateOfBirth  Date   `json:"dob,omitempty"`
	MobileCode   int    `json:"mobile_code,omitempty"`
	MobileNumber int    `json:"mobile_number,omitempty"`
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
