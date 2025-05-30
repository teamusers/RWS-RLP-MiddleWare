// services.go
package services

import (
	"context"
	"fmt"
	"strings"

	"lbe/api/http/responses"
	"lbe/config"
	"lbe/model"
	"lbe/utils"
	"net/http"
)

const (
	// Endpoints
	CreateProfileURL = "/priv/v1/apps/:api_key/users"
	ProfileURL       = "/priv/v1/apps/:api_key/external/users"
	EventUrl         = "/api/1.0/user_events/:event_name"

	// Event Names
	RlpEventNameUpdateUserTier = "update_user_tier"
)

func CreateProfile(ctx context.Context, client *http.Client, payload any) (*responses.GetUserResponse, []byte, error) {
	return profile(ctx, client, http.MethodPost, buildRlpProfileURL(CreateProfileURL, "", ""), payload)
}

func UpdateProfile(ctx context.Context, client *http.Client, externalId string, payload any) (*responses.GetUserResponse, []byte, error) {
	return profile(ctx, client, http.MethodPut, buildRlpProfileURL(ProfileURL, externalId, ""), payload)
}

func GetProfile(ctx context.Context, client *http.Client, externalId string) (*responses.GetUserResponse, []byte, error) {
	query := "user[user_profile]=true&expand_incentives=true&show_identifiers=true"
	return profile(ctx, client, http.MethodGet, buildRlpProfileURL(ProfileURL, externalId, query), nil)
}

// TODO: update logic
func UpdateUserTier(ctx context.Context, client *http.Client, payload any) (*responses.UserTierUpdateEventResponse, []byte, error) {
	conf := config.GetConfig()
	endpoint := strings.ReplaceAll(EventUrl, ":event_name", RlpEventNameUpdateUserTier)
	urlWithParams := fmt.Sprintf("%s%s", conf.Api.Rlp.Host, endpoint)

	return utils.DoAPIRequest[responses.UserTierUpdateEventResponse](model.APIRequestOptions{
		Method:         http.MethodPost,
		URL:            urlWithParams,
		Body:           payload,
		ExpectedStatus: http.StatusOK,
		Client:         client,
		Context:        ctx,
		ContentType:    model.ContentTypeJson,
	})
}

func profile(ctx context.Context, client *http.Client, operation, url string, payload any) (*responses.GetUserResponse, []byte, error) {
	conf := config.GetConfig()

	return utils.DoAPIRequest[responses.GetUserResponse](model.APIRequestOptions{
		Method: operation,
		URL:    url,
		Body:   payload,
		BasicAuth: &model.BasicAuthCredentials{
			Username: conf.Api.Rlp.Core.ApiKey,
			Password: conf.Api.Rlp.Core.ApiSecret,
		},
		ExpectedStatus: http.StatusOK,
		Client:         client,
		Context:        ctx,
		ContentType:    model.ContentTypeJson,
	})
}

func buildRlpProfileURL(basePath, externalId, queryParams string) string {
	conf := config.GetConfig()
	endpoint := strings.ReplaceAll(basePath, ":api_key", conf.Api.Rlp.Core.ApiKey)

	if externalId != "" {
		endpoint = fmt.Sprintf("%s/%s", endpoint, externalId)
	}

	if queryParams != "" {
		endpoint = fmt.Sprintf("%s?%s", endpoint, queryParams)
	}

	return fmt.Sprintf("%s%s", conf.Api.Rlp.Core.Host, endpoint)
}
