// services.go
package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"lbe/api/http/responses"
	"lbe/config"
	"net/http"
	"time"
)

// Endpoints
const (
	authURL  = "/api/v1/auth"
	usersURL = "/api/v1/user/login"
)

var ErrRecordNotFound = errors.New("record not found")

func GetAccessToken() (string, error) {
	// Create the request body according to the given specification.

	AppID := config.GetConfig().API.Memberservice.AppID
	secretKey := config.GetConfig().API.Memberservice.Secret
	reqBody, err := GenerateSignature(AppID, secretKey)

	if err != nil {
		// Handle the error appropriately.
		log.Fatalf("unable to generate auth signature: %v", err)
	}

	// Encode the request body into JSON.
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Create a new HTTP request.
	req, err := http.NewRequest("GET", BuildFullURL(authURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Set required headers.
	req.Header.Set("AppID", AppID)
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client with a timeout.
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for a non-OK status.
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("authentication endpoint returned non-OK status")
	}

	// Decode the response.
	var authResp responses.MemberAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", err
	}

	return authResp.Data.AccessToken, nil
}

// GetUserByEmail first gets an access token, then calls the users endpoint using the token
// to query a user by email. It returns a UserResponse or an error.
func GetLoginUserByEmail(email string) (*responses.UserResponse, error) {
	// Get the access token.
	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}

	urlWithEmail := fmt.Sprintf("%s/%s", BuildFullURL(usersURL), email)
	req, err := http.NewRequest("GET", urlWithEmail, nil)
	if err != nil {
		return nil, err
	}

	// Set the Bearer token and required headers.
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AppID", "app123")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// For demonstration, we assume that if the email is not found,
	// the endpoint returns HTTP Status 404. If so, we simulate a gorm.ErrRecordNotFound.
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrRecordNotFound
	}

	// For any other non-OK status, return an error.
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(token)
	}

	// Decode the response.
	var userResp responses.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, err
	}

	return &userResp, nil
}
