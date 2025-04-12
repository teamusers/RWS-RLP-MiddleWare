package requests

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}
