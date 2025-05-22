package responses

const (
	// error messages
	CiamUserAlreadyExists = "Another object with the same value for property userPrincipalName already exists."
)

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

// GraphUserCollection holds the “value” array from a Graph `/users` response.
type GraphUserCollection struct {
	Value []GraphUser `json:"value"`
}

type GraphUser struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
}

// GraphCreateUserResponse extracts the CIAM EEID user id of the newly created user from the response.
type GraphCreateUserResponse struct {
	Id string `json:"id"`
	// other fields are irrelevant
}

// GraphUserExtensionCollection holds the “value” array from a Graph `/users` response.
type GraphUserExtensionCollection struct {
	Value []struct {
		ID                string `json:"id"`
		DisplayName       string `json:"displayName"`
		Mail              string `json:"mail"`
		UserPrincipalName string `json:"userPrincipalName"`
	} `json:"value"`
}

type GraphApiErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	// other fields ignored
}
