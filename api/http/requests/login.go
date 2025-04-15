package requests

type Login struct {
	Email string `json:"email" binding:"required"`
}
