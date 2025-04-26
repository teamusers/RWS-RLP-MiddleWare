// services.go
package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/config"
	"lbe/model"
	"net/http"
	"time"
)

// Endpoints
const (
	authURL                = "/api/v1/auth"
	userURL                = "/api/v1/user"
	verifyUserExistenceURL = "/api/v1/user/verify"
	registerURL            = "/api/v1/user/register"
	loginURL               = "/api/v1/user/login"
	updateBurnPinURL       = "/api/v1/user/pin"
)

func GetAccessToken() (string, error) {
	AppID := config.GetConfig().Api.Memberservice.AppID
	secretKey := config.GetConfig().Api.Memberservice.Secret
	reqBody, err := GenerateSignature(AppID, secretKey)

	if err != nil {
		log.Printf("unable to generate auth signature: %v", err)
		return "", err
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
	var authResp responses.ApiResponse[responses.MemberAuthResponseData]
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", err
	}

	return authResp.Data.AccessToken, nil
}

// requests a check to verify member existence via email,
// login session token is returned if updateSessionToken query param is true
func VerifyMemberExistence(email string, updateSessionToken bool) (*responses.ApiResponse[model.LoginSessionToken], error) {
	payload := requests.VerifyUser{
		Email: email,
	}

	var targetURL string
	if updateSessionToken {
		targetURL = loginURL
	} else {
		targetURL = verifyUserExistenceURL
	}

	resp, err := buildMemberHttpClient("POST", BuildFullURL(targetURL), payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userResp responses.ApiResponse[model.LoginSessionToken]
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, err
	}

	return &userResp, nil
}

// TODO: update accordingly when member sevice endpoint updates
func PostRegisterUser(payload requests.CreateUser) error {

	_, err := buildMemberHttpClient("POST", BuildFullURL(registerURL), payload)
	if err != nil {
		return err
	}
	return nil

}

func UpdateBurnPin(payload interface{}) error {
	// Build the full URL by combining the base URL, email, and signUpType.
	urlWithParams := BuildFullURL(updateBurnPinURL)

	resp, err := buildMemberHttpClient("PUT", urlWithParams, payload)
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

func buildMemberHttpClient(httpMethod string, url string, payload any) (*http.Response, error) {
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

	fmt.Print(url)
	// Create a new POST request with the JSON payload.
	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set the required headers.
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AppID", config.GetConfig().Api.Memberservice.AppID)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}
