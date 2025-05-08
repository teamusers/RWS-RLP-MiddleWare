package v1

import (
	"net/http"

	"lbe/api/http/requests"
	"lbe/api/http/responses"

	"lbe/api/http/services"

	"github.com/gin-gonic/gin"
)

// VerifyUserHandler handles POST /ciam/users/verify
func VerifyUserHandler(c *gin.Context) {
	var req requests.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coll, err := services.VerifyCIAMExistence(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(coll.Value) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no user found"})
		return
	}

	c.JSON(http.StatusOK, coll)
}

// RegisterUserHandler handles POST /ciam/users/register
func RegisterUserHandler(c *gin.Context) {
	var req requests.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1) Map slice of requests.Identity → []responses.Identity
	identities := make([]responses.Identity, len(req.Identities))
	for i, id := range req.Identities {
		identities[i] = responses.Identity{
			SignInType:       id.SignInType,
			Issuer:           id.Issuer,
			IssuerAssignedID: id.IssuerAssignedID,
		}
	}

	// 2) Build the GraphCreateUserPayload
	payload := responses.GraphCreateUserPayload{
		AccountEnabled:   req.AccountEnabled,
		DisplayName:      req.DisplayName,
		MailNickname:     req.MailNickname,
		Identities:       identities,
		Mail:             req.Mail,
		PasswordPolicies: req.PasswordPolicies,
	}

	// Manually map into the anonymous struct
	payload.PasswordProfile = struct {
		ForceChangePasswordNextSignIn bool   `json:"forceChangePasswordNextSignIn"`
		Password                      string `json:"password"`
	}{
		ForceChangePasswordNextSignIn: req.PasswordProfile.ForceChangePasswordNextSignIn,
		Password:                      req.PasswordProfile.Password,
	}

	if err := services.PostCIAMRegisterUser(c.Request.Context(), payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered"})
}
