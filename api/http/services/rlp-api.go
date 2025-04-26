// services.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"lbe/api/http/responses"
	"lbe/config"
	"net/http"
)

// Endpoints
const (
	UpdateProfileURL = "/priv/v1/apps/:api_key/users"
	ProfileURL       = "/priv/v1/apps/:api_key/external/users/:external_id"
)

func Profile(external_id string, payload any, operation string, endpoint string) (*responses.GetUserResponse, error) {
	conf := config.GetConfig()
	endpoint = strings.ReplaceAll(endpoint, ":api_key", conf.Api.Rlp.ApiKey)
	endpoint = strings.ReplaceAll(endpoint, ":external_id", external_id)
	urlWithParams := fmt.Sprintf("%s%s", conf.Api.Rlp.Host, endpoint)

	resp, err := buildHttpClient(operation, urlWithParams, payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// In case of error reading the body, just return the status code
			return nil, fmt.Errorf("error calling RLP services: received status code %d, but failed to read response body: %v", resp.StatusCode, err)
		}
		defer resp.Body.Close()

		return nil, fmt.Errorf("error calling RLP services: received status code %d and response: %s", resp.StatusCode, string(body))
	}

	// on 200 OK, read & clean the body
	raw, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 1) strip UTF-8 BOM if present
	raw = bytes.TrimPrefix(raw, []byte("\xef\xbb\xbf"))
	// 2) replace any non-breaking space (U+00A0) with a normal space
	raw = []byte(strings.ReplaceAll(string(raw), "\u00A0", " "))

	// now unmarshal into your struct
	var result responses.GetUserResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w; cleaned body: %q", err, raw)
	}

	return &result, nil
}
