package requests

import "lbe/model"

type Register struct {
	Email string `json:"email" binding:"required"`
}

type RegisterGr struct {
	GrId  string `json:"gr_id" binding:"required"`
	GrPin string `json:"gr_pin" binding:"required"`
}

type RegisterGrCms struct {
	model.GrMember
	Url string `json:"url" binding:"required"`
}
