package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"lbe/api/http/controllers/v1/user"
	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/codes"
	"lbe/config"
	"lbe/model"
	"lbe/system"
	"lbe/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestGrTierMatching(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		expectErr bool
	}{
		{
			name:      "Tier A - GR class 1",
			input:     "Class 1",
			expected:  "",
			expectErr: false,
		},
		{
			name:      "Tier B - GR class 2",
			input:     "Class 2",
			expected:  "Tier B",
			expectErr: false,
		},
		{
			name:      "Tier C - GR class 3",
			input:     "Class 3",
			expected:  "Tier C",
			expectErr: false,
		},
		{
			name:      "Tier C - GR class 4",
			input:     "Class 4",
			expected:  "Tier C",
			expectErr: false,
		},
		{
			name:      "Tier C - GR class 5",
			input:     "Class 5",
			expected:  "Tier C",
			expectErr: false,
		},
		{
			name:      "Tier D - GR class 6",
			input:     "Class 6",
			expected:  "Tier D",
			expectErr: false,
		},
		{
			name:      "Invalid format - only one part",
			input:     "Class",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Invalid format - three parts",
			input:     "Class 1 extra",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Invalid format - non-integer level",
			input:     "Class X",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Invalid format - class level < 1",
			input:     "Class 0",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "Invalid format - negative level",
			input:     "Class -2",
			expected:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := user.GrTierMatching(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestVerifyUserExistence(t *testing.T) {
	defer gock.Off()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/user/register/verify", user.VerifyUserExistence)

	tests := []struct {
		name                    string
		requestBody             any
		setupMocks              func(email string)
		expectedHTTPCode        int
		expectedResponseCode    int64
		expectedResponseMessage string
	}{
		{
			name:        "User not found - success",
			requestBody: requests.VerifyUserExistence{Email: "newuser@example.com"},
			setupMocks: func(email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user not found
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildEmailFilter(email)).
					Reply(200).
					JSON(map[string]any{"value": []any{}})

				// Mock ACS auth
				gock.New(config.GetConfig().Api.Acs.Host).
					Post(services.AcsAuthURL).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock ACS email send
				sendEndpoint := strings.ReplaceAll(
					services.AcsSendEmailByTemplateURL,
					":template_name", services.AcsEmailTemplateRequestOtp,
				)
				gock.New(config.GetConfig().Api.Acs.Host).
					Post(sendEndpoint).
					Reply(200).
					JSON(nil)
			},
			expectedHTTPCode:        http.StatusOK,
			expectedResponseCode:    codes.SUCCESSFUL,
			expectedResponseMessage: "existing user not found",
		},
		{
			name:        "User already exists - conflict",
			requestBody: requests.VerifyUserExistence{Email: "existing@example.com"},
			setupMocks: func(email string) {
				// CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// CIAM user exists
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildEmailFilter(email)).
					Reply(200).
					JSON(map[string]any{
						"value": []map[string]any{{"id": "abc123"}},
					})
			},
			expectedHTTPCode:        http.StatusConflict,
			expectedResponseCode:    codes.EXISTING_USER_FOUND,
			expectedResponseMessage: "existing user found",
		},
		{
			name:        "CIAM get user by email fail - error",
			requestBody: requests.VerifyUserExistence{Email: "existing@example.com"},
			setupMocks: func(email string) {
				// CIAM auth error
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(500)
			},
			expectedHTTPCode:        http.StatusInternalServerError,
			expectedResponseCode:    codes.INTERNAL_ERROR,
			expectedResponseMessage: "internal error",
		},
		{
			name:        "ACS send email fail - error",
			requestBody: requests.VerifyUserExistence{Email: "newuser@example.com"},
			setupMocks: func(email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user not found
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildEmailFilter(email)).
					Reply(200).
					JSON(map[string]any{"value": []any{}})

				// Mock ACS auth error
				gock.New(config.GetConfig().Api.Acs.Host).
					Post(services.AcsAuthURL).
					Reply(500)
			},
			expectedHTTPCode:        http.StatusInternalServerError,
			expectedResponseCode:    codes.INTERNAL_ERROR,
			expectedResponseMessage: "internal error",
		},

		{
			name:                    "Invalid request body - error",
			requestBody:             nil, // raw string invalid body
			setupMocks:              func(email string) {},
			expectedHTTPCode:        http.StatusBadRequest,
			expectedResponseCode:    codes.INVALID_REQUEST_BODY,
			expectedResponseMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			// Setup mocks
			if req, ok := tt.requestBody.(requests.VerifyUserExistence); ok {
				tt.setupMocks(req.Email)
			}

			var bodyBytes []byte
			switch b := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(b)
			default:
				bodyBytes, _ = json.Marshal(b)
			}

			req := httptest.NewRequest(http.MethodPost, "/user/register/verify", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedHTTPCode, rec.Code)

			if rec.Code == http.StatusOK || rec.Code == http.StatusConflict {
				var resp responses.ApiResponse[model.Otp]
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponseCode, resp.Code)
				assert.Equal(t, tt.expectedResponseMessage, resp.Message)
			}
		})
	}
}

// LBE 6 Unit Test
func TestVerifyGrExistence(t *testing.T) {
	defer gock.Off()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/user/gr", user.VerifyGrExistence)

	validSampleReq, _ := utils.LoadTestData[requests.VerifyGrUser]("lbe6_verifyGrExistence_req.json")

	tests := []struct {
		name                    string
		requestBody             any
		setupMocks              func(grId, email string)
		expectedHTTPCode        int
		expectedResponseCode    int64
		expectedResponseMessage string
	}{
		{
			name:        "GR ID not found - success",
			requestBody: validSampleReq,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user by GR ID returns empty
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildGrIdFilter(grId)).
					Reply(200).
					JSON(responses.GraphUserCollection{})

				// CMS member fetch
				gock.New(config.GetConfig().Api.Cms.Host).
					Get(services.GetMemberURL).
					MatchParam("systemId", config.GetConfig().Api.Cms.SystemID).
					MatchParam("memberId", grId).
					Reply(200).
					JSON(responses.GRProfilePayload{
						EmailAddress:       email,
						ContactOptionEmail: true,
					})

				// Mock ACS auth
				gock.New(config.GetConfig().Api.Acs.Host).
					Post(services.AcsAuthURL).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock ACS email send
				sendEndpoint := strings.ReplaceAll(
					services.AcsSendEmailByTemplateURL,
					":template_name", services.AcsEmailTemplateRequestOtp,
				)
				gock.New(config.GetConfig().Api.Acs.Host).
					Post(sendEndpoint).
					Reply(200).
					JSON(nil)
			},
			expectedHTTPCode:        http.StatusOK,
			expectedResponseCode:    codes.SUCCESSFUL,
			expectedResponseMessage: "gr profile found",
		},
		{
			name:        "GR ID already linked - conflict",
			requestBody: validSampleReq,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// CIAM returns existing user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildGrIdFilter(grId)).
					Reply(200).
					JSON(responses.GraphUserCollection{
						Value: []responses.GraphUser{{ID: "existing-user"}},
					})
			},
			expectedHTTPCode:        http.StatusConflict,
			expectedResponseCode:    codes.GR_MEMBER_LINKED,
			expectedResponseMessage: "gr profile already linked to another email",
		},
		{
			name:        "CIAM get user by grId fail - error",
			requestBody: validSampleReq,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(500)
			},
			expectedHTTPCode:        http.StatusInternalServerError,
			expectedResponseCode:    codes.INTERNAL_ERROR,
			expectedResponseMessage: "internal error",
		},
		{
			name:        "CMS profile fetch - error",
			requestBody: validSampleReq,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user by GR ID returns empty
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildGrIdFilter(grId)).
					Reply(200).
					JSON(responses.GraphUserCollection{})

				// CMS member fetch error
				gock.New(config.GetConfig().Api.Cms.Host).
					Get(services.GetMemberURL).
					MatchParam("systemId", config.GetConfig().Api.Cms.SystemID).
					MatchParam("memberId", grId).
					Reply(500)
			},
			expectedHTTPCode:        http.StatusInternalServerError,
			expectedResponseCode:    codes.INTERNAL_ERROR,
			expectedResponseMessage: "internal error",
		},
		{
			name:        "ACS send email fail - error",
			requestBody: validSampleReq,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user by GR ID returns empty
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildGrIdFilter(grId)).
					Reply(200).
					JSON(responses.GraphUserCollection{})

				// CMS member fetch
				gock.New(config.GetConfig().Api.Cms.Host).
					Get(services.GetMemberURL).
					MatchParam("systemId", config.GetConfig().Api.Cms.SystemID).
					MatchParam("memberId", grId).
					Reply(200).
					JSON(responses.GRProfilePayload{
						EmailAddress:       email,
						ContactOptionEmail: true,
					})

				// Mock ACS auth error
				gock.New(config.GetConfig().Api.Acs.Host).
					Post(services.AcsAuthURL).
					Reply(500)
			},
			expectedHTTPCode:        http.StatusInternalServerError,
			expectedResponseCode:    codes.INTERNAL_ERROR,
			expectedResponseMessage: "internal error",
		},
		{
			name:                    "Invalid request body - error",
			requestBody:             `{}`,
			setupMocks:              func(email, grId string) {}, // No mock needed
			expectedHTTPCode:        http.StatusBadRequest,
			expectedResponseCode:    codes.INVALID_REQUEST_BODY,
			expectedResponseMessage: "",
		},
		{
			name:                    "Invalid JSON ShouldBindJSON - error",
			requestBody:             `{"user": "invalid-json}`,   // malformed JSON (missing closing quote)
			setupMocks:              func(email, grId string) {}, // No mocks needed
			expectedHTTPCode:        http.StatusBadRequest,
			expectedResponseCode:    codes.INVALID_REQUEST_BODY,
			expectedResponseMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			var grId, email string
			if req, ok := tt.requestBody.(requests.VerifyGrUser); ok {
				grId = req.User.GrProfile.Id
				email = req.User.Email
			}

			tt.setupMocks(grId, email)

			var bodyBytes []byte
			switch b := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(b)
			default:
				bodyBytes, _ = json.Marshal(b)
			}

			req := httptest.NewRequest(http.MethodPost, "/user/gr", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedHTTPCode, rec.Code)

			if rec.Code == http.StatusOK || rec.Code == http.StatusConflict {
				var resp responses.ApiResponse[any]
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponseCode, resp.Code)
				assert.Equal(t, tt.expectedResponseMessage, resp.Message)
			}
		})
	}
}

func TestVerifyGrCmsExistence(t *testing.T) {
	defer gock.Off()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/user/gr-cms", user.VerifyGrCmsExistence)

	validSampleReq, _ := utils.LoadTestData[requests.VerifyGrCmsUser]("lbe7_verifyGrCmsExistence_req.json")

	tests := []struct {
		name                    string
		requestBody             any
		setupMocks              func(email, grId string)
		expectedHTTPCode        int
		expectedResponseCode    int64
		expectedResponseMessage string
	}{
		{
			name:        "User and GR ID not found - success",
			requestBody: validSampleReq,
			setupMocks: func(email, grId string) {

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user by email returns empty
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildEmailFilter(email)).
					Reply(200).
					JSON(responses.GraphUserCollection{})

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user by GR ID returns empty
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildGrIdFilter(grId)).
					Reply(200).
					JSON(responses.GraphUserCollection{})

					// Mock ACS auth
				gock.New(config.GetConfig().Api.Acs.Host).
					Post(services.AcsAuthURL).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock ACS email send //TODO: update to correct template
				sendEndpoint := strings.ReplaceAll(
					services.AcsSendEmailByTemplateURL,
					":template_name", services.AcsEmailTemplateRequestOtp,
				)
				gock.New(config.GetConfig().Api.Acs.Host).
					Post(sendEndpoint).
					Reply(200).
					JSON(nil)
			},
			expectedHTTPCode:        http.StatusOK,
			expectedResponseCode:    codes.SUCCESSFUL,
			expectedResponseMessage: "existing user not found",
		},
		{
			name:        "Existing email found - conflict",
			requestBody: validSampleReq,
			setupMocks: func(email, grId string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user by email returning found user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildEmailFilter(email)).
					Reply(200).
					JSON(responses.GraphUserCollection{
						Value: []responses.GraphUser{
							{ID: "abc123"},
						},
					})
			},
			expectedHTTPCode:        http.StatusConflict,
			expectedResponseCode:    codes.EXISTING_USER_FOUND,
			expectedResponseMessage: "existing user found",
		},
		{
			name:        "Existing GR ID found - conflict",
			requestBody: validSampleReq,
			setupMocks: func(email, grId string) {

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user by email returns empty
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildEmailFilter(email)).
					Reply(200).
					JSON(responses.GraphUserCollection{})

				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user by GR ID returns a found user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildGrIdFilter(grId)).
					Reply(200).
					JSON(responses.GraphUserCollection{
						Value: []responses.GraphUser{
							{ID: "abc123"},
						},
					})
			},
			expectedHTTPCode:        http.StatusConflict,
			expectedResponseCode:    codes.GR_MEMBER_LINKED,
			expectedResponseMessage: "gr profile already linked to another email",
		},
		{
			name:        "CIAM get user by email fail - error",
			requestBody: validSampleReq,
			setupMocks: func(email, grId string) {

				// Mock CIAM auth error
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(500)
			},
			expectedHTTPCode:        http.StatusInternalServerError,
			expectedResponseCode:    codes.INTERNAL_ERROR,
			expectedResponseMessage: "internal error",
		},
		{
			name:        "CIAM get user by grId fail - error",
			requestBody: validSampleReq,
			setupMocks: func(email, grId string) {

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user by email returns empty
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildEmailFilter(email)).
					Reply(200).
					JSON(responses.GraphUserCollection{})

				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user by GR ID error
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildGrIdFilter(grId)).
					Reply(500)
			},
			expectedHTTPCode:        http.StatusInternalServerError,
			expectedResponseCode:    codes.INTERNAL_ERROR,
			expectedResponseMessage: "internal error",
		},
		{
			name:        "ACS send email fail - error",
			requestBody: validSampleReq,
			setupMocks: func(email, grId string) {

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user by email returns empty
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildEmailFilter(email)).
					Reply(200).
					JSON(responses.GraphUserCollection{})

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user by GR ID returns empty
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", buildGrIdFilter(grId)).
					Reply(200).
					JSON(responses.GraphUserCollection{})

					// Mock ACS auth error
				gock.New(config.GetConfig().Api.Acs.Host).
					Post(services.AcsAuthURL).
					Reply(500)
			},
			expectedHTTPCode:        http.StatusInternalServerError,
			expectedResponseCode:    codes.INTERNAL_ERROR,
			expectedResponseMessage: "internal error",
		},
		{
			name:                    "Invalid request body - error",
			requestBody:             nil, // raw string invalid body
			setupMocks:              func(email, grId string) {},
			expectedHTTPCode:        http.StatusBadRequest,
			expectedResponseCode:    0,
			expectedResponseMessage: "",
		},
		{
			name:                    "Invalid JSON ShouldBindJSON - error",
			requestBody:             `{"user": "invalid-json}`,   // malformed JSON (missing closing quote)
			setupMocks:              func(email, grId string) {}, // No mocks needed
			expectedHTTPCode:        http.StatusBadRequest,
			expectedResponseCode:    0,
			expectedResponseMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			var email, grId string
			if req, ok := tt.requestBody.(requests.VerifyGrCmsUser); ok {
				email = req.User.Email
				if req.User.GrProfile != nil {
					grId = req.User.GrProfile.Id
				}
			}

			tt.setupMocks(email, grId)

			var bodyBytes []byte
			switch b := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(b)
			default:
				bodyBytes, _ = json.Marshal(b)
			}

			req := httptest.NewRequest(http.MethodPost, "/user/gr-cms", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedHTTPCode, rec.Code)

			if rec.Code == http.StatusOK || rec.Code == http.StatusConflict {
				var resp responses.ApiResponse[any]
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponseCode, resp.Code)
				assert.Equal(t, tt.expectedResponseMessage, resp.Message)
			}
		})
	}
}

func TestGetCachedGrCmsProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                    string
		regId                   string
		mockCacheHit            bool
		mockCacheValue          *model.User
		mockCacheErr            error
		expectedHTTPCode        int
		expectedResponseCode    int64
		expectedResponseMessage string
	}{
		{
			name:                    "Cache miss returns conflict",
			regId:                   "123",
			mockCacheHit:            false,
			mockCacheErr:            errors.New("not found"),
			expectedHTTPCode:        http.StatusConflict,
			expectedResponseCode:    codes.CACHED_PROFILE_NOT_FOUND,
			expectedResponseMessage: "cached profile not found",
		},
		{
			name:                    "Cache hit returns success",
			regId:                   "123",
			mockCacheHit:            true,
			mockCacheValue:          &model.User{DateOfBirth: model.GetDatePointer("2000-01-01")},
			expectedHTTPCode:        http.StatusOK,
			expectedResponseCode:    codes.SUCCESSFUL,
			expectedResponseMessage: "cached profile found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//setup cache
			if tt.mockCacheValue != nil {
				system.ObjectSet(tt.regId, tt.mockCacheValue, 30*time.Minute)
				defer system.ObjectDelete(tt.regId)
			} else {
				system.ObjectDelete(tt.regId) // ensure clean start if no mock value
			}

			router := gin.New()
			router.GET("/gr-cms/:reg_id", user.GetCachedGrCmsProfile)

			req := httptest.NewRequest(http.MethodGet, "/gr-cms/"+tt.regId, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedHTTPCode, rec.Code)

			if tt.expectedHTTPCode == http.StatusOK || tt.expectedHTTPCode == http.StatusConflict {
				var resp responses.ApiResponse[responses.VerifyGrCmsUserResponseData]
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				if err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				assert.Equal(t, tt.expectedResponseCode, resp.Code)
				assert.Equal(t, tt.expectedResponseMessage, resp.Message)
				if tt.mockCacheHit {
					assert.Equal(t, tt.regId, resp.Data.RegId)
					assert.Equal(t, *tt.mockCacheValue.DateOfBirth, resp.Data.DateOfBirth)
				}
			}
		})
	}
}

// test utils

func buildEmailFilter(email string) string {
	return fmt.Sprintf("mail eq '%s'", email)
}

func buildGrIdFilter(grId string) string {
	return fmt.Sprintf("%s/grid eq '%s'", config.GetConfig().Api.Eeid.UserIdLinkExtensionKey, grId)
}
