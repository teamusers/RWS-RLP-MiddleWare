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
	ProfileURL = "/priv/v1/apps/:api_key/external/users"
	EventUrl   = "/api/1.0/user_events/:event_name"

	// Event Names
	RlpEventNameUpdateUserTier = "update_user_tier"
)

func PutProfile(ctx context.Context, client *http.Client, externalId string, payload any) (*responses.GetUserResponse, error) {
	return profile(ctx, client, externalId, payload, http.MethodPut)
}

func GetProfile(ctx context.Context, client *http.Client, externalId string) (*responses.GetUserResponse, error) {
	return profile(ctx, client, externalId, nil, http.MethodGet)
}

func UpdateUserTier(ctx context.Context, client *http.Client, payload any) (*responses.UserTierUpdateEventResponse, error) {
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

func profile(ctx context.Context, client *http.Client, externalId string, payload any, operation string) (*responses.GetUserResponse, error) {
	conf := config.GetConfig()
	endpoint := strings.ReplaceAll(ProfileURL, ":api_key", conf.Api.Rlp.ApiKey)
	if externalId != "" {
		endpoint = fmt.Sprintf("%s/%s", endpoint, externalId)
	}
	urlWithParams := fmt.Sprintf("%s%s", conf.Api.Rlp.Host, endpoint)

	return utils.DoAPIRequest[responses.GetUserResponse](model.APIRequestOptions{
		Method:         operation,
		URL:            urlWithParams,
		Body:           payload,
		ExpectedStatus: http.StatusOK,
		Client:         client,
		Context:        ctx,
		ContentType:    model.ContentTypeJson,
	})
}
