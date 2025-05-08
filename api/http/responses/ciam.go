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

// Identity represents one sign-in identity for a Graph user.
type Identity struct {
	SignInType       string `json:"signInType"`       // e.g. "emailAddress"
	Issuer           string `json:"issuer"`           // e.g. "eeidtest1.onmicrosoft.com"
	IssuerAssignedID string `json:"issuerAssignedId"` // e.g. "ue.test1@eeidtest1.onmicrosoft.com"
}

// GraphCreateUserPayload is the full shape to create a new Azure AD (EEID) user.
type GraphCreateUserPayload struct {
	AccountEnabled  bool       `json:"accountEnabled"`
	DisplayName     string     `json:"displayName"`
	MailNickname    string     `json:"mailNickname"`
	Identities      []Identity `json:"identities"`
	Mail            string     `json:"mail"`
	PasswordProfile struct {
		ForceChangePasswordNextSignIn bool   `json:"forceChangePasswordNextSignIn"`
		Password                      string `json:"password"`
	} `json:"passwordProfile"`
	PasswordPolicies string `json:"passwordPolicies"`
}
