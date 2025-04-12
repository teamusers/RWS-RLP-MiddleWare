package requests

type AuthRequest struct {
	Timestamp string `json:"Timestamp"`
	Nonce     string `json:"Nonce"`
	Signature string `json:"Signature"`
}
