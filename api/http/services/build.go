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
	var req *http.Request
	var err error

	client := &http.Client{Timeout: 10 * time.Second}

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

	resp, err := client.Do(req)
	return resp, err
}

// BuildFullURL constructs the full endpoint URL using the host from the configuration
// and appending the provided endpoint.
func BuildFullURL(endpoint string) string {
	conf := config.GetConfig() // Get the centralized configuration.
	host := conf.Api.Memberservice.Host
	return fmt.Sprintf("%s%s", host, endpoint)
}
