package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lbe/api/http/responses"
	"lbe/config"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	//endpoints
	AcsAuthURL                = "/api/v1/auth"
	AcsSendEmailByTemplateURL = "/api/v1/send/template/:template_name"

	// subjects
	AcsEmailSubjectRequestOtp = "RWS Loyalty Program - Verify OTP"

	// template names
	AcsEmailTemplateRequestOtp = "request_email_otp"
)

func getAcsAccessToken() (string, error) {
	appId := config.GetConfig().Api.Acs.AppId
	secretKey := config.GetConfig().Api.Acs.Secret
	reqBody, err := GenerateSignature(appId, secretKey)

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
	req, err := http.NewRequest("POST", buildFullAcsUrl(AcsAuthURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("AppID", appId)
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
	var authResp responses.ApiResponse[responses.AcsAuthResponseData]
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", err
	}

	return authResp.Data.AccessToken, nil
}

func PostAcsSendEmailByTemplate(templateName string, payload any) error {
	url := strings.ReplaceAll(AcsSendEmailByTemplateURL, ":template_name", templateName)

	if _, err := buildAcsHttpClient(http.MethodPost, buildFullAcsUrl(url), payload, http.StatusOK); err != nil {
		return err
	}

	return nil
}

func buildAcsHttpClient(httpMethod, url string, payload any, expectedStatus int) (*http.Response, error) {
	var req *http.Request
	var err error

	// Get the access token.
	token, err := getAcsAccessToken()
	if err != nil {
		return nil, err
	}

	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("error marshaling payload: %w", err)
		}
		req, _ = http.NewRequest(httpMethod, url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(httpMethod, url, nil)
	}
	if err != nil {
		return nil, err
	}

	// Set the required headers.
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("AppID", config.GetConfig().Api.Acs.AppId)
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

	return resp, err
}

func buildFullAcsUrl(endpoint string) string {
	return fmt.Sprintf("%s%s", config.GetConfig().Api.Acs.Host, endpoint)
}
