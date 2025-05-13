package responses

type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

// GraphUserCollection holds the “value” array from a Graph `/users` response.
type GraphUserCollection struct {
	Value []struct {
		ID                string `json:"id"`
		DisplayName       string `json:"displayName"`
		Mail              string `json:"mail"`
		UserPrincipalName string `json:"userPrincipalName"`
	} `json:"value"`
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
