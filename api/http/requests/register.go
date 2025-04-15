package requests

type Register struct {
	Email string `json:"email" binding:"required"`
}
