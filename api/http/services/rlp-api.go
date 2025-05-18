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

// Endpoints
const (
	ProfileURL = "/priv/v1/apps/:api_key/external/users/:external_id"
)

func PutProfile(ctx context.Context, client *http.Client, external_id string, payload any) (*responses.GetUserResponse, error) {
	return profile(ctx, client, external_id, payload, http.MethodPut)
}

func GetProfile(ctx context.Context, client *http.Client, external_id string) (*responses.GetUserResponse, error) {
	return profile(ctx, client, external_id, nil, http.MethodGet)
}

func profile(ctx context.Context, client *http.Client, external_id string, payload any, operation string) (*responses.GetUserResponse, error) {
	conf := config.GetConfig()
	endpoint := strings.ReplaceAll(ProfileURL, ":api_key", conf.Api.Rlp.ApiKey)
	endpoint = strings.ReplaceAll(endpoint, ":external_id", external_id)
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
