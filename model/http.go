package model

import (
	"context"
	"net/http"
)

type ctxKey string

const (
	// content type
	ContentTypeJson = "application/json"
	ContentTypeForm = "application/x-www-form-urlencoded"

	// context keys
	HttpClientCtxKey ctxKey = "httpClient"
)

type APIRequestOptions struct {
	Method         string
	URL            string
	Body           any
	BearerToken    string
	ExpectedStatus int
	Headers        map[string]string
	Client         *http.Client
	Context        context.Context
	ContentType    string
}
