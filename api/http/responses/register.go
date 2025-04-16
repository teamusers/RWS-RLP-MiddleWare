package responses

import "lbe/model"

type GetGrMemberResponse struct {
	model.GrMember
	model.Otp
}
