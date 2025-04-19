package responses

// APIResponse is the standard envelope for successful operations.
// The Data field contains the payload, which varies by endpoint.
type APIResponse struct {
	// Message provides a human‑readable status or result description.
	// Example: "user created", "email found"
	Message string `json:"message" example:""`

	// Data holds the response payload. Its type depends on the endpoint:
	// e.g. AuthResponse for /auth, LoginResponse for /user/login, etc.
	Data interface{} `json:"data"`
}

// ErrorResponse is used when reporting simple error messages.
type ErrorResponse struct {
	// Error provides the error detail.
	// Example: "invalid json request body"
	Error string `json:"error" example:"invalid json request body"`
}
