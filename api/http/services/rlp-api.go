// services.go
package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"lbe/api/http/responses"
	"lbe/config"
	"net/http"
)

// Endpoints
const (
	MemberProfileURL    = "/priv/v1/apps/:api_key/external/users/:external_id"
	NewMemberURL        = "/priv/v1/apps/:api_key/users"
	TierUpdate          = "/priv/v1/apps/:api_key/external/users/:external_id/events"
	UpdateMemberProfile = "/priv/v1/apps/:api_key/external/users/:external_id"
)

var ErrRecordNotFound = errors.New("record not found")

func Member(external_id string, payload any, operation string, endpoint string) (*responses.GetRlpMemberUserResponse, error) {
	conf := config.GetConfig()
	endpoint = strings.ReplaceAll(endpoint, ":api_key", conf.Api.Rlp.AppID)
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

	// Decode the response.
	var Resp responses.GetRlpMemberUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&Resp); err != nil {
		return nil, err
	}

	return &Resp, nil
}
