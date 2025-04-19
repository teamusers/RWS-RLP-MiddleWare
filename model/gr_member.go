package model

// GrMember represents a GR CMS member’s profile information.
//
// Example:
// {
//   "gr_id": "GR12345",
//   "f_name": "Jane",
//   "l_name": "Doe",
//   "email": "jane.doe@example.com",
//   "dob": "1985-04-12",
//   "mobile": "98765432"
// }
type GrMember struct {
	// GrId is the unique GR member identifier.
	GrId *string `json:"gr_id" example:"GR12345"`

	// FirstName is the given name of the GR member.
	FirstName *string `json:"f_name" example:"Jane"`

	// LastName is the family name of the GR member.
	LastName *string `json:"l_name" example:"Doe"`

	// Email is the member's email address.
	Email *string `json:"email" example:"jane.doe@example.com"`

	// Dob is the date of birth in YYYY-MM-DD format.
	Dob *string `json:"dob" example:"1985-04-12"`

	// Mobile is the contact phone number.
	Mobile *string `json:"mobile" example:"98765432"`
}
