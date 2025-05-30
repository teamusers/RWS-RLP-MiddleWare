package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lbe/model"
	"log"
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

func DoAPIRequest[T any](opts model.APIRequestOptions) (*T, []byte, error) {
	var bodyReader io.Reader
	if opts.Body != nil {
		switch opts.ContentType {
		case model.ContentTypeForm:
			// Expect Body to be url.Values
			if form, ok := opts.Body.(url.Values); ok {
				bodyReader = strings.NewReader(form.Encode())
			} else {
				return nil, nil, fmt.Errorf("body must be url.Values for content type %s", opts.ContentType)
			}
		case model.ContentTypeJson:
			fallthrough
		default:
			// Marshal as JSON
			b, err := json.Marshal(opts.Body)
			if err != nil {
				return nil, nil, fmt.Errorf("marshaling body: %w", err)
			}
			bodyReader = bytes.NewReader(b)
			opts.ContentType = model.ContentTypeJson // ensure default if empty
		}
	}

	req, err := http.NewRequestWithContext(opts.Context, opts.Method, opts.URL, bodyReader)
	if err != nil {
		return nil, nil, fmt.Errorf("creating request: %w", err)
	}

	if opts.Body != nil {
		req.Header.Set("Content-Type", opts.ContentType)
	}

	// add only one type of auth
	if opts.BearerToken != "" && opts.BasicAuth != nil {
		return nil, nil, errors.New("cannot use both Bearer token and Basic Auth")
	}

	if opts.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+opts.BearerToken)
	}

	if opts.BasicAuth != nil {
		req.SetBasicAuth(opts.BasicAuth.Username, opts.BasicAuth.Password)
	}

	// Add any other headers
	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	// --- LOG REQUEST HERE ---
	log.Printf("[API REQUEST] %s %s; Content-Type: %s; Body: %s", opts.Method, opts.URL, opts.ContentType, func() string {
		if opts.Body == nil {
			return "<empty>"
		}
		b, err := json.Marshal(opts.Body)
		if err != nil {
			return "<error marshaling body>"
		}
		return string(b)
	}())

	resp, err := opts.Client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("reading response body: %w", err)
	}

	// --- LOG RESPONSE HERE ---
	log.Printf("[API RESPONSE] Status: %d; Body: %s", resp.StatusCode,
		strings.Join(strings.Fields(strings.ReplaceAll(strings.ReplaceAll(string(raw), "\n", " "), "\t", " ")), " "))

	// 1) strip UTF-8 BOM if present
	raw = bytes.TrimPrefix(raw, []byte("\xef\xbb\xbf"))
	// 2) replace any non-breaking space (U+00A0) with a normal space
	raw = []byte(strings.ReplaceAll(string(raw), "\u00A0", " "))

	if resp.StatusCode != opts.ExpectedStatus {
		return nil, raw, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(raw))
	}

	if len(raw) == 0 {
		var empty T
		return &empty, raw, nil
	}

	var result T
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, raw, fmt.Errorf("failed to decode response: %w; cleaned body: %q", err, raw)
	}
	return &result, raw, nil
}
