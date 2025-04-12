package requests

type AuthRequest struct {
	Timestamp *string `json:"timestamp" binding:"required"`
	Nonce     *string `json:"nonce" binding:"required"`
	Sign      *string `json:"sign" binding:"required"`
}

type InitiateLoginRequest struct {
	Email     *string `json:"email" binding:"required,email"`
	DebugMode *int    `json:"debugMode"` //to be removed!
}
