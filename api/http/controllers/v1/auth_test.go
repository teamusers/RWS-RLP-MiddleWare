package v1_test

import (
	"bytes"
	"encoding/json"
	v1 "lbe/api/http/controllers/v1"
	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/codes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func Test_AuthHandler(t *testing.T) {
	defer gock.Off()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth", v1.AuthHandler)

	// TODO: refactor logic and remove hardcode
	appId := "app1234"
	timestamp := "1747979126"
	nonce := "API"
	signature := "f02ab7387e10629136ea15e1e7a2537f93a407f5c204b8efa39d491a76ab1faa"

	tests := []struct {
		name                    string
		appID                   string
		requestBody             any
		setupMocks              func(appID string)
		expectedHTTPCode        int
		expectedResponseCode    int64
		expectedResponseMessage string
		expectedResponseBody    any
	}{
		{
			name:  "SUCCESS - valid signature",
			appID: appId,
			requestBody: requests.AuthRequest{
				Nonce:     nonce,
				Timestamp: timestamp,
				Signature: signature,
			},
			setupMocks: func(appID string) {
			},
			expectedHTTPCode:        http.StatusOK,
			expectedResponseCode:    codes.SUCCESSFUL,
			expectedResponseMessage: "token successfully generated",
		},
		{
			name:  "UNAUTHORIZED - missing app id",
			appID: "",
			requestBody: requests.AuthRequest{
				Nonce:     nonce,
				Timestamp: timestamp,
				Signature: signature,
			},
			setupMocks: func(appID string) {
			},
			expectedHTTPCode:     http.StatusUnauthorized,
			expectedResponseBody: responses.MissingAppIdErrorResponse(),
		},
		{
			name:  "UNAUTHORIZED - invalid app id",
			appID: "1",
			requestBody: requests.AuthRequest{
				Nonce:     nonce,
				Timestamp: timestamp,
				Signature: signature,
			},
			setupMocks: func(appID string) {
			},
			expectedHTTPCode:     http.StatusUnauthorized,
			expectedResponseBody: responses.InvalidAppIdErrorResponse(),
		},
		{
			name:  "UNAUTHORIZED - invalid signature",
			appID: appId,
			requestBody: requests.AuthRequest{
				Nonce:     nonce,
				Timestamp: timestamp,
				Signature: "123",
			},
			setupMocks: func(appID string) {
			},
			expectedHTTPCode:     http.StatusUnauthorized,
			expectedResponseBody: responses.InvalidSignatureErrorResponse(),
		},
		{
			name:  "ERROR - invalid req body",
			appID: appId,
			requestBody: requests.AuthRequest{
				Nonce: nonce,
			},
			setupMocks: func(appID string) {
			},
			expectedHTTPCode:     http.StatusBadRequest,
			expectedResponseBody: responses.InvalidRequestBodyErrorResponse(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			tt.setupMocks(tt.appID)

			var bodyBytes []byte
			switch b := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(b)
			default:
				bodyBytes, _ = json.Marshal(b)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			if tt.appID != "" {
				req.Header.Set("AppID", tt.appID)
			}
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedHTTPCode, rec.Code)

			if rec.Code == http.StatusOK {
				var resp responses.ApiResponse[any]
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedResponseCode, resp.Code)
				assert.Equal(t, tt.expectedResponseMessage, resp.Message)
			} else if tt.expectedResponseBody != nil {
				expected, _ := json.Marshal(tt.expectedResponseBody)
				actual := rec.Body.Bytes()
				assert.JSONEq(t, string(expected), string(actual))
			}
		})
	}
}

func Test_InvalidQueryParametersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/test", v1.InvalidQueryParametersHandler)

	tests := []struct {
		name                 string
		expectedHTTPCode     int
		expectedResponseBody any
	}{
		{
			name:                 "ERROR - invalid / missing query param",
			expectedHTTPCode:     http.StatusBadRequest,
			expectedResponseBody: responses.InvalidQueryParametersErrorResponse(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedHTTPCode, rec.Code)

			var resp responses.ApiResponse[any]
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)

			expected, _ := json.Marshal(tt.expectedResponseBody)
			actual, _ := json.Marshal(resp)
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}
