// services.go
package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"lbe/api/http/responses"
	"lbe/codes"
	"lbe/config"
	"net/http"
	"time"
)

// Endpoints
const (
	authURL          = "/api/v1/auth"
	userURL = "/api/v1/user"
	updateBurnPinURL = "/api/v1/user/pin"
)

var ErrRecordNotFound = errors.New("record not found")

func GetAccessToken() (string, error) {
	AppID := config.GetConfig().API.Memberservice.AppID
	secretKey := config.GetConfig().API.Memberservice.Secret
	reqBody, err := GenerateSignature(AppID, secretKey)

	if err != nil {
		log.Fatalf("unable to generate auth signature: %v", err)
	}

	// Encode the request body into JSON.
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Create a new HTTP request.
	req, err := http.NewRequest("POST", BuildFullURL(authURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

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
// to query a user by email. It returns a Login session token or an error.
func GetLoginUserByEmail(email string) (*responses.MemberLoginResponse, error) {
	urlWithEmail := fmt.Sprintf("%s?updateSessionToken=true", BuildFullURL(userURL))

	payload := map[string]string{
		"email": email,
	}

	resp, err := buildHttpClient("POST", urlWithEmail, payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
		case http.StatusOK:
			var userResp responses.MemberLoginResponse
			if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
				return nil, err
			}
			return &userResp, nil
		case 201:
			return nil, ErrRecordNotFound
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("error calling member services: status %d, response: %s", resp.StatusCode, string(body))
	}
}


func GetRegisterUserByEmail(email string) error {
	urlWithParams := BuildFullURL(userURL)
	payload := map[string]string{
		"email": email,
	}

	resp, err := buildHttpClient("POST", urlWithParams, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		return nil 
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("error calling member services: status %d, response: %s", resp.StatusCode, string(body))
}


// PostRegisterUser posts a JSON payload to register a user by combining email and signUpType in the URL.
func PostRegisterUser(payload interface{}) error {

	// Build the full URL by combining the base URL, email, and signUpType.
	urlWithParams := BuildFullURL(userURL)

	resp, err := buildHttpClient("POST", urlWithParams, payload)
	if err != nil {
		return err
	}

	// Check for a non-OK status.
	switch resp.StatusCode {
	case http.StatusOK:
		return fmt.Errorf("error calling member services: received status code %d and response: %s", 201, "User already exists")
	case 201:
		return nil
	default:
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error calling member services: status %d, response: %s", resp.StatusCode, string(body))
}

}

func UpdateBurnPin(payload interface{}) error {
	// Build the full URL by combining the base URL, email, and signUpType.
	urlWithParams := BuildFullURL(updateBurnPinURL)

	resp, err := buildHttpClient("PUT", urlWithParams, payload)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// In case of error reading the body, just return the status code
			return fmt.Errorf("error calling member services: received status code %d, but failed to read response body: %v", resp.StatusCode, err)
		}
		defer resp.Body.Close()

		return fmt.Errorf("error calling member services: received status code %d and response: %s", resp.StatusCode, string(body))
	}

	return nil
}

func buildHttpClient(httpMethod string, url string, payload any) (*http.Response, error) {
	// Get the access token.
	token, err := GetAccessToken()
	if err != nil {
		return nil, err
	}

	// Marshal the passed payload into JSON.
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	// Create a new POST request with the JSON payload.
	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set the required headers.
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AppID", config.GetConfig().API.Memberservice.AppID)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}
