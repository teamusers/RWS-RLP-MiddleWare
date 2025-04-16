package model

type GrMember struct {
	GrId      *string `json:"gr_id"`
	FirstName *string `json:"f_name"`
	LastName  *string `json:"l_name"`
	Email     *string `json:"email"`
	Dob       *string `json:"dob"`
	Mobile    *string `json:"mobile"`
}
