package requests

type AuthRequest struct {
	Timestamp string `json:"timestamp" binding:"required"`
	Nonce     string `json:"nonce" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}
