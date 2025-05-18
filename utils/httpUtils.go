package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lbe/model"
	"net/http"
	"net/url"
	"strings"
)

func WithHttpClient(ctx context.Context, client *http.Client) context.Context {
	return context.WithValue(ctx, model.HttpClientCtxKey, client)
}

func GetHttpClient(ctx context.Context) *http.Client {
	if client, ok := ctx.Value(model.HttpClientCtxKey).(*http.Client); ok {
		return client
	}
	return http.DefaultClient // fallback
}

func DoAPIRequest[T any](opts model.APIRequestOptions) (*T, error) {
	var bodyReader io.Reader
	if opts.Body != nil {
		switch opts.ContentType {
		case model.ContentTypeForm:
			// Expect Body to be url.Values
			if form, ok := opts.Body.(url.Values); ok {
				bodyReader = strings.NewReader(form.Encode())
			} else {
				return nil, fmt.Errorf("body must be url.Values for content type %s", opts.ContentType)
			}
		case model.ContentTypeJson:
			fallthrough
		default:
			// Marshal as JSON
			b, err := json.Marshal(opts.Body)
			if err != nil {
				return nil, fmt.Errorf("marshaling body: %w", err)
			}
			bodyReader = bytes.NewReader(b)
			opts.ContentType = model.ContentTypeJson // ensure default if empty
		}
	}

	req, err := http.NewRequestWithContext(opts.Context, opts.Method, opts.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if opts.Body != nil {
		req.Header.Set("Content-Type", opts.ContentType)
	}

	if opts.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+opts.BearerToken)
	}

	// Add any other headers
	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	resp, err := opts.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != opts.ExpectedStatus {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	if len(respBody) == 0 {
		var empty T
		return &empty, nil
	}

	var result T
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result, nil
}
