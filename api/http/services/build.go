package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"lbe/config"
)

func buildHttpClient(httpMethod string, url string, payload any) (*http.Response, error) {

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

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

// BuildFullURL constructs the full endpoint URL using the host from the configuration
// and appending the provided endpoint.
func BuildFullURL(endpoint string) string {
	conf := config.GetConfig() // Get the centralized configuration.
	host := conf.Api.Memberservice.Host
	return fmt.Sprintf("%s%s", host, endpoint)
}
