// services.go
package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/config"
	"lbe/model"
	"lbe/utils"
)

// Endpoints
const (
	CiamAuthURL      = "/oauth2/v2.0/token"
	CiamUserURL      = "/v1.0/users"
	extensionURL     = "/v1.0/users/:id/extensions"
	extensionDataURL = "/v1.0/users/:id/extensions/:extensionsid"
)

// GetCIAMAccessToken acquires a bearer token from Azure AD using client credentials.
func GetCIAMAccessToken(ctx context.Context, client *http.Client) (*responses.TokenResponse, []byte, error) {
	cfg := config.GetConfig().Api.Eeid

	host := strings.TrimRight(cfg.AuthHost, "/")
	tenantID := cfg.TenantID
	clientID := cfg.ClientID
	clientSecret := cfg.ClientSecret

	tokenURL := fmt.Sprintf("%s/%s%s", host, tenantID, CiamAuthURL)

	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"scope":         {"https://graph.microsoft.com/.default"},
	}

	return utils.DoAPIRequest[responses.TokenResponse](model.APIRequestOptions{
		Method:         http.MethodPost,
		URL:            tokenURL,
		Body:           form,
		ExpectedStatus: http.StatusOK,
		Client:         client,
		Context:        ctx,
		ContentType:    model.ContentTypeForm,
	})
}

// GetCIAMUserByEmail calls Graph GET /users?$filter=mail eq '{email}'.
func GetCIAMUserByEmail(ctx context.Context, client *http.Client, email string) (*responses.GraphUserCollection, []byte, error) {
	tokenResp, _, err := GetCIAMAccessToken(ctx, client)
	if err != nil {
		return nil, nil, fmt.Errorf("getting access token: %w", err)
	}
	// extract the actual bearer token
	bearer := tokenResp.AccessToken

	cfg := config.GetConfig().Api.Eeid
	base := strings.TrimRight(cfg.Host, "/")
	filter := url.QueryEscape(fmt.Sprintf("mail eq '%s'", email))
	fullURL := fmt.Sprintf("%s%s?$filter=%s", base, CiamUserURL, filter)

	return utils.DoAPIRequest[responses.GraphUserCollection](model.APIRequestOptions{
		Method:         http.MethodGet,
		URL:            fullURL,
		Body:           nil,
		BearerToken:    bearer,
		ExpectedStatus: http.StatusOK,
		Client:         client,
		Context:        ctx,
		ContentType:    model.ContentTypeJson,
	})
}

// GetCIAMUserByGrId calls Graph GET /users?$filter={schemaIdKey}/grid eq '{grId}'.
func GetCIAMUserByGrId(ctx context.Context, client *http.Client, grId string) (*responses.GraphUserCollection, []byte, error) {
	tokenResp, _, err := GetCIAMAccessToken(ctx, client)
	if err != nil {
		return nil, nil, fmt.Errorf("getting access token: %w", err)
	}
	// extract the actual bearer token
	bearer := tokenResp.AccessToken

	cfg := config.GetConfig().Api.Eeid
	base := strings.TrimRight(cfg.Host, "/")
	filter := url.QueryEscape(fmt.Sprintf("%s/grid eq '%s'", cfg.UserIdLinkExtensionKey, grId))
	fullURL := fmt.Sprintf("%s%s?$filter=%s", base, CiamUserURL, filter)

	return utils.DoAPIRequest[responses.GraphUserCollection](model.APIRequestOptions{
		Method:         http.MethodGet,
		URL:            fullURL,
		Body:           nil,
		BearerToken:    bearer,
		ExpectedStatus: http.StatusOK,
		Client:         client,
		Context:        ctx,
		ContentType:    model.ContentTypeJson,
	})
}

// PostCIAMRegisterUser calls Graph POST /users to create a new AD user.
func PostCIAMRegisterUser(ctx context.Context, client *http.Client, payload requests.GraphCreateUserRequest) (*responses.GraphCreateUserResponse, []byte, error) {
	tokenResp, _, err := GetCIAMAccessToken(ctx, client)
	if err != nil {
		return nil, nil, fmt.Errorf("getting access token: %w", err)
	}
	// extract the actual bearer token
	bearer := tokenResp.AccessToken

	cfg := config.GetConfig().Api.Eeid
	base := strings.TrimRight(cfg.Host, "/")
	fullURL := fmt.Sprintf("%s%s", base, CiamUserURL)

	log.Printf("registering user: %v", payload.Mail)
	return utils.DoAPIRequest[responses.GraphCreateUserResponse](model.APIRequestOptions{
		Method:         http.MethodPost,
		URL:            fullURL,
		Body:           payload,
		BearerToken:    bearer,
		ExpectedStatus: http.StatusCreated,
		Client:         client,
		Context:        ctx,
		ContentType:    model.ContentTypeJson,
	})
}

func PatchCIAMAddUserSchemaExtensions(ctx context.Context, client *http.Client, userId string, payload any) ([]byte, error) {
	tokenResp, _, err := GetCIAMAccessToken(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("getting access token: %w", err)
	}
	// extract the actual bearer token
	bearer := tokenResp.AccessToken

	cfg := config.GetConfig().Api.Eeid
	base := strings.TrimRight(cfg.Host, "/")
	fullURL := fmt.Sprintf("%s%s/%s", base, CiamUserURL, userId)

	_, raw, err := utils.DoAPIRequest[struct{}](model.APIRequestOptions{
		Method:         http.MethodPatch,
		URL:            fullURL,
		Body:           payload,
		BearerToken:    bearer,
		ExpectedStatus: http.StatusNoContent,
		Client:         client,
		Context:        ctx,
		ContentType:    model.ContentTypeJson,
	})
	return raw, err
}

func PatchCIAMUpdateUser(ctx context.Context, client *http.Client, userId string, payload any) ([]byte, error) {
	tokenResp, _, err := GetCIAMAccessToken(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("getting access token: %w", err)
	}
	// extract the actual bearer token
	bearer := tokenResp.AccessToken

	cfg := config.GetConfig().Api.Eeid
	base := strings.TrimRight(cfg.Host, "/")
	fullURL := fmt.Sprintf("%s%s/%s", base, CiamUserURL, userId)

	log.Printf("patching CIAM user id: %v", userId)
	_, raw, err := utils.DoAPIRequest[struct{}](model.APIRequestOptions{
		Method:         http.MethodPatch,
		URL:            fullURL,
		Body:           payload,
		BearerToken:    bearer,
		ExpectedStatus: http.StatusNoContent,
		Client:         client,
		Context:        ctx,
		ContentType:    model.ContentTypeJson,
	})
	return raw, err
}
