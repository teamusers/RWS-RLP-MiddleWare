package requests

type AuthRequest struct {
	// Unix timestamp (seconds since epoch) when the request was generated.
	Timestamp string `json:"timestamp" binding:"required" example:"1744075148"`

	// A unique random string for each request to prevent replay attacks.
	Nonce string `json:"nonce" binding:"required" example:"API"`

	// HMAC-SHA256 signature of "appID|timestamp|nonce" hex-encoded.
	// Computed by concatenating the appID, timestamp, and nonce to form a base string,
	// then applying HMAC-SHA256 with the secret key and hex-encoding the resulting digest.
	Signature string `json:"signature" binding:"required" example:"1558850cb1b48e826197c48d6a14c5f3bf4b644bcb0065ceb0b07978296116bc"`
}
