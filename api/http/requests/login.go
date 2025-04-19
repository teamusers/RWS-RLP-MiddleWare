package requests

type Login struct {
	// Email address of the user attempting to log in
	Email string `json:"email" binding:"required" example:"user@example.com"`
}
