// requests/register.go
package requests

// VerifyRequest is the payload for verifying a user by email.
type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Identity represents one sign-in identity for a user.
type Identity struct {
	SignInType       string `json:"signInType" binding:"required,oneof=emailAddress userName"`
	Issuer           string `json:"issuer" binding:"required"`
	IssuerAssignedID string `json:"issuerAssignedId" binding:"required"`
}

// PasswordProfile holds the password settings for a new user.
type PasswordProfile struct {
	ForceChangePasswordNextSignIn bool   `json:"forceChangePasswordNextSignIn"`
	Password                      string `json:"password" binding:"required,min=8"`
}

// RegisterRequest is the full payload for creating a new EEID user.
type RegisterRequest struct {
	AccountEnabled   bool            `json:"accountEnabled"`
	DisplayName      string          `json:"displayName" binding:"required"`
	MailNickname     string          `json:"mailNickname" binding:"required"`
	Identities       []Identity      `json:"identities" binding:"required,dive,required"`
	Mail             string          `json:"mail" binding:"required,email"`
	PasswordProfile  PasswordProfile `json:"passwordProfile"`
	PasswordPolicies string          `json:"passwordPolicies"`
}
