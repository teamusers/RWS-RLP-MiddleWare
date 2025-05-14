// services.go
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/config"
)

// Endpoints
const (
	ciamAuthURL      = "/oauth2/v2.0/token"
	ciamUserURL      = "/v1.0/users"
	extensionURL     = "/v1.0/users/:id/extensions"
	extensionDataURL = "/v1.0/users/:id/extensions/:extensionsid"
)

// GetCIAMAccessToken acquires a bearer token from Azure AD using client credentials.
func GetCIAMAccessToken(ctx context.Context) (*responses.TokenResponse, error) {
	cfg := config.GetConfig().Api.Eeid

	host := strings.TrimRight(cfg.AuthHost, "/")
	tenantID := cfg.TenantID
	clientID := cfg.ClientID
	clientSecret := cfg.ClientSecret

	tokenURL := fmt.Sprintf("%s/%s%s", host, tenantID, ciamAuthURL)

	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"scope":         {"https://graph.microsoft.com/.default"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("CIAM token endpoint error %d: %s", resp.StatusCode, body)
		return nil, fmt.Errorf("failed to acquire token, status %d", resp.StatusCode)
	}

	var tr responses.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, err
	}
	return &tr, nil
}

// GetCIAMUserByEmail calls Graph GET /users?$filter=mail eq '{email}'.
func GetCIAMUserByEmail(ctx context.Context, email string) (*responses.GraphUserCollection, error) {
	tokenResp, err := GetCIAMAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting access token: %w", err)
	}
	// extract the actual bearer token
	bearer := tokenResp.AccessToken

	cfg := config.GetConfig().Api.Eeid
	base := strings.TrimRight(cfg.Host, "/")
	filter := url.QueryEscape(fmt.Sprintf("mail eq '%s'", email))
	fullURL := fmt.Sprintf("%s%s?$filter=%s", base, ciamUserURL, filter)

	return doJSONRequest[responses.GraphUserCollection](ctx, http.MethodGet, fullURL, bearer, nil, http.StatusOK)
}

// GetCIAMUserByGrId calls Graph GET /users?$filter={schemaIdKey}/grid eq '{grId}'.
func GetCIAMUserByGrId(ctx context.Context, grId string) (*responses.GraphUserCollection, error) {
	tokenResp, err := GetCIAMAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting access token: %w", err)
	}
	// extract the actual bearer token
	bearer := tokenResp.AccessToken

	cfg := config.GetConfig().Api.Eeid
	base := strings.TrimRight(cfg.Host, "/")
	filter := url.QueryEscape(fmt.Sprintf("%s/grid eq '%s'", cfg.UserIdLinkExtensionKey, grId))
	fullURL := fmt.Sprintf("%s%s?$filter=%s", base, ciamUserURL, filter)

	return doJSONRequest[responses.GraphUserCollection](ctx, http.MethodGet, fullURL, bearer, nil, http.StatusOK)
}

// PostCIAMRegisterUser calls Graph POST /users to create a new AD user.
func PostCIAMRegisterUser(ctx context.Context, payload requests.GraphCreateUserRequest) (*responses.GraphCreateUserResponse, error) {
	tokenResp, err := GetCIAMAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting access token: %w", err)
	}
	// extract the actual bearer token
	bearer := tokenResp.AccessToken

	cfg := config.GetConfig().Api.Eeid
	base := strings.TrimRight(cfg.Host, "/")
	fullURL := fmt.Sprintf("%s%s", base, ciamUserURL)

	log.Printf("registering user: %v", payload.DisplayName)
	return doJSONRequest[responses.GraphCreateUserResponse](ctx, http.MethodPost, fullURL, bearer, payload, http.StatusCreated)
}

func PatchCIAMAddUserSchemaExtensions(ctx context.Context, userId string, payload any) error {
	tokenResp, err := GetCIAMAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("getting access token: %w", err)
	}
	// extract the actual bearer token
	bearer := tokenResp.AccessToken

	cfg := config.GetConfig().Api.Eeid
	base := strings.TrimRight(cfg.Host, "/")
	fullURL := fmt.Sprintf("%s%s/%s", base, ciamUserURL, userId)

	_, err = doJSONRequest[struct{}](ctx, http.MethodPatch, fullURL, bearer, payload, http.StatusNoContent)
	return err
}

func PatchCIAMUpdateUser(ctx context.Context, userId string, payload any) error {
	tokenResp, err := GetCIAMAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("getting access token: %w", err)
	}
	// extract the actual bearer token
	bearer := tokenResp.AccessToken

	cfg := config.GetConfig().Api.Eeid
	base := strings.TrimRight(cfg.Host, "/")
	fullURL := fmt.Sprintf("%s%s/%s", base, ciamUserURL, userId)

	log.Printf("patching CIAM user id: %v", userId)
	_, err = doJSONRequest[struct{}](ctx, http.MethodPatch, fullURL, bearer, payload, http.StatusNoContent)
	return err
}

func doJSONRequest[T any](ctx context.Context, method, url string, bearerToken string, body any, expectedStatus int) (*T, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	if len(respBody) == 0 {
		// Empty body, return nil or a zero-value
		var empty T
		return &empty, nil
	}

	var result T
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result, nil
}
