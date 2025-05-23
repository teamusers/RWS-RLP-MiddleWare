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
	GetMemberURL = "/cms-webapi-bsp/v2/member"
)

// TODO: Fix to correct spec
func GRMemberProfile(memberId string, payload any, operation string, endpoint string) (*responses.GRProfilePayload, error) {
	conf := config.GetConfig()
	urlWithParams := fmt.Sprintf("%s%s?systemId=%s&memberId=%s", conf.Api.Cms.Host, endpoint, conf.Api.Cms.SystemID, memberId)

	resp, err := buildHttpClient(operation, urlWithParams, payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// In case of error reading the body, just return the status code
			return nil, fmt.Errorf("error calling CMS services: received status code %d, but failed to read response body: %v", resp.StatusCode, err)
		}
		defer resp.Body.Close()

		return nil, fmt.Errorf("error calling CMS services: received status code %d and response: %s", resp.StatusCode, string(body))
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
	var result responses.GRProfilePayload
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w; cleaned body: %q", err, raw)
	}

	return &result, nil
}
