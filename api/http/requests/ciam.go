// requests/register.go
package requests

import (
	"fmt"
	"lbe/model"
	"strings"
)

// VerifyRequest is the payload for verifying a user by email.
type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Identity represents one sign-in identity for a Graph user.
type Identity struct {
	SignInType       string `json:"signInType"`       // e.g. "emailAddress"
	Issuer           string `json:"issuer"`           // e.g. "eeidtest1.onmicrosoft.com"
	IssuerAssignedID string `json:"issuerAssignedId"` // e.g. "ue.test1@eeidtest1.onmicrosoft.com"
}

// PasswordProfile holds the password settings for a new user.
type PasswordProfile struct {
	ForceChangePasswordNextSignIn bool   `json:"forceChangePasswordNextSignIn"`
	Password                      string `json:"password" binding:"required,min=8"`
}

// GraphCreateUserRequest is the full shape payload to create a new Azure AD (EEID) user.
type GraphCreateUserRequest struct {
	AccountEnabled   bool            `json:"accountEnabled"`
	DisplayName      string          `json:"displayName"`
	MailNickname     string          `json:"mailNickname"`
	Identities       []Identity      `json:"identities"`
	Mail             string          `json:"mail"`
	PasswordProfile  PasswordProfile `json:"passwordProfile"`
	PasswordPolicies string          `json:"passwordPolicies"`
}

func GenerateInitialRegistrationRequest(user *model.User) GraphCreateUserRequest {
	return GraphCreateUserRequest{
		AccountEnabled: true,
		DisplayName:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		MailNickname:   user.Email[:strings.Index(user.Email, "@")],
		Identities: []Identity{
			{
				SignInType:       "emailAddress",
				Issuer:           user.Email[strings.Index(user.Email, "@")+1:],
				IssuerAssignedID: user.Email,
			},
		},
		Mail: user.Email,
		PasswordProfile: PasswordProfile{
			ForceChangePasswordNextSignIn: true,
			Password:                      "P@ssw0rd!2025",
		},
		PasswordPolicies: "DisablePasswordExpiration",
	}
}

type UserIdLinkSchemaExtensionFields struct {
	RlpId string `json:"rlpid"`
	RlpNo string `json:"rlpno"`
	GrId  string `json:"grid"`
}

type GraphDisableAccountRequest struct {
	AccountEnabled bool `json:"accountEnabled"`
}
