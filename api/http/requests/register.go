package requests

type Register struct {
	Email      string `json:"email"`
	SignUpType string `json:"sign_up_type"`
}
