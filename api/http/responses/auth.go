package responses

// AuthResponse is returned upon successful authentication.
// swagger:response AuthResponse
type AuthResponseData struct {
	// AccessToken is the JWT issued to the client for subsequent requests.
	// Example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}
