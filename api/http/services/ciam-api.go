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

// VerifyCIAMExistence calls Graph GET /users?$filter=mail eq '{email}'.
func VerifyCIAMExistence(ctx context.Context, email string) (*responses.GraphUserCollection, error) {
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}
	// extract the actual bearer token
	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Graph /users error %d: %s", resp.StatusCode, body)
		return nil, fmt.Errorf("failed to verify user existence; status %d", resp.StatusCode)
	}

	var coll responses.GraphUserCollection
	if err := json.NewDecoder(resp.Body).Decode(&coll); err != nil {
		return nil, err
	}
	return &coll, nil
}

// PostCIAMRegisterUser calls Graph POST /users to create a new AD user.
func PostCIAMRegisterUser(ctx context.Context, payload responses.GraphCreateUserPayload) error {
	tokenResp, err := GetCIAMAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("getting access token: %w", err)
	}
	// extract the actual bearer token
	bearer := tokenResp.AccessToken

	cfg := config.GetConfig().Api.Eeid
	base := strings.TrimRight(cfg.Host, "/")
	fullURL := fmt.Sprintf("%s%s", base, ciamUserURL)

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+bearer)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to register user; status %d: %s", resp.StatusCode, body)
	}
	return nil
}
