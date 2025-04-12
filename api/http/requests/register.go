package requests

type SignUpRequest struct {
	Email      string `json:"email" binding:"required,email"`
	SignUpType string `json:"sign_up_type" binding:"required"`
}
