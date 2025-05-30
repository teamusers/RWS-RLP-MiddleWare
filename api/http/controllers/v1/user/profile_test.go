package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lbe/api/http/controllers/v1/user"
	"lbe/api/http/requests"
	"lbe/api/http/responses"
	"lbe/api/http/services"
	"lbe/config"
	"lbe/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

// LBE 9 Unit Test
func Test_LBE_9_GetUserProfile(t *testing.T) {
	defer gock.Off()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/user/:external_id", user.GetUserProfile)

	rlpGetProfileRes := utils.LoadTestData[responses.GetUserResponse]("rlp_update_profile_res.json")
	rlpUserProfileNotFoundRes := utils.LoadTestData[responses.UserProfileErrorResponse]("rlp_user_not_found_err_res.json")
	expectedRes := utils.LoadTestData[responses.ApiResponse[any]]("lbe9_getUser_res.json")

	rlpProfileUrl := strings.ReplaceAll(services.ProfileURL, ":api_key", config.GetConfig().Api.Rlp.Core.ApiKey)

	tests := []struct {
		name                 string
		externalId           string
		requestBody          any
		setupMocks           func()
		expectedHTTPCode     int
		expectedResponseBody any
	}{
		{
			name:       "SUCCESS - User found",
			externalId: "25052300047",
			setupMocks: func() {
				// Mock RLP Get Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Get(rlpProfileUrl).
					Reply(200).
					JSON(rlpGetProfileRes)
			},
			expectedHTTPCode:     http.StatusOK,
			expectedResponseBody: expectedRes,
		},
		{
			name:       "CONFLICT - User not found",
			externalId: "25052300047",
			setupMocks: func() {
				// Mock RLP Get Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Get(rlpProfileUrl).
					Reply(400).
					JSON(rlpUserProfileNotFoundRes)
			},
			expectedHTTPCode:     http.StatusConflict,
			expectedResponseBody: responses.ExistingUserNotFoundErrorResponse(),
		},
		{
			name:       "ERROR - RLP Get user fail",
			externalId: "25052300047",
			setupMocks: func() {
				// Mock RLP Get Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Get(rlpProfileUrl).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			// Setup mocks
			tt.setupMocks()

			var bodyBytes []byte
			switch b := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(b)
			default:
				bodyBytes, _ = json.Marshal(b)
			}

			req := httptest.NewRequest(http.MethodGet, "/user/"+tt.externalId, bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedHTTPCode, rec.Code)

			if rec.Code == http.StatusOK || rec.Code == http.StatusConflict {
				var resp responses.ApiResponse[any]
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)

				if tt.expectedResponseBody != nil {
					expected, _ := json.Marshal(tt.expectedResponseBody)
					actual, _ := json.Marshal(resp)
					assert.JSONEq(t, string(expected), string(actual))
				}
			}
		})
	}
}

// LBE 10 Unit Test
func Test_LBE_10_UpdateUserProfile(t *testing.T) {
	defer gock.Off()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/user/update/:external_id", user.UpdateUserProfile)

	rlpUpdateProfileRes := utils.LoadTestData[responses.GetUserResponse]("rlp_update_profile_update_res.json")
	rlpUserProfileNotFoundRes := utils.LoadTestData[responses.UserProfileErrorResponse]("rlp_user_not_found_err_res.json")

	validSampleReq := utils.LoadTestData[requests.UpdateUserProfile]("lbe10_updateUser_req.json")
	expectedRes := utils.LoadTestData[responses.ApiResponse[any]]("lbe10_updateUser_res.json")

	rlpProfileUrl := strings.ReplaceAll(services.ProfileURL, ":api_key", config.GetConfig().Api.Rlp.Core.ApiKey)

	tests := []struct {
		name                 string
		externalId           string
		requestBody          any
		setupMocks           func()
		expectedHTTPCode     int
		expectedResponseBody any
	}{
		{
			name:        "SUCCESS - User updated",
			externalId:  "25052300047",
			requestBody: validSampleReq,
			setupMocks: func() {
				// Mock RLP Put Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(rlpProfileUrl).
					Reply(200).
					JSON(rlpUpdateProfileRes)
			},
			expectedHTTPCode:     http.StatusOK,
			expectedResponseBody: expectedRes,
		},
		{
			name:        "CONFLICT - User not found",
			externalId:  "25052300047",
			requestBody: validSampleReq,
			setupMocks: func() {
				// Mock RLP Put Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(rlpProfileUrl).
					Reply(400).
					JSON(rlpUserProfileNotFoundRes)
			},
			expectedHTTPCode:     http.StatusConflict,
			expectedResponseBody: responses.ExistingUserNotFoundErrorResponse(),
		},
		{
			name:        "ERROR - RLP put user fail",
			externalId:  "25052300047",
			requestBody: validSampleReq,
			setupMocks: func() {
				// Mock RLP Put Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(rlpProfileUrl).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:                 "ERROR - Invalid JSON ShouldBindJSON",
			externalId:           "25052300047",
			requestBody:          `{"user": "invalid-json}`, // malformed JSON (missing closing quote)
			setupMocks:           func() {},                 // No mocks needed
			expectedHTTPCode:     http.StatusBadRequest,
			expectedResponseBody: responses.InvalidRequestBodyErrorResponse(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			// Setup mocks
			tt.setupMocks()

			var bodyBytes []byte
			switch b := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(b)
			default:
				bodyBytes, _ = json.Marshal(b)
			}

			req := httptest.NewRequest(http.MethodPut, "/user/update/"+tt.externalId, bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedHTTPCode, rec.Code)

			if rec.Code == http.StatusOK || rec.Code == http.StatusConflict {
				var resp responses.ApiResponse[any]
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)

				if tt.expectedResponseBody != nil {
					expected, _ := json.Marshal(tt.expectedResponseBody)
					actual, _ := json.Marshal(resp)
					assert.JSONEq(t, string(expected), string(actual))
				}
			}
		})
	}
}

// LBE 11 Unit Test
func Test_LBE_11_WithdrawUserProfile(t *testing.T) {
	defer gock.Off()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PUT("/user/archive/:external_id", user.WithdrawUserProfile)

	rlpGetProfileRes := utils.LoadTestData[responses.GetUserResponse]("rlp_update_profile_res.json")
	rlpUpdateProfileRes := utils.LoadTestData[responses.GetUserResponse]("rlp_update_profile_withdraw_res.json")
	rlpUserProfileNotFoundRes := utils.LoadTestData[responses.UserProfileErrorResponse]("rlp_user_not_found_err_res.json")

	ciamGetAuth := utils.LoadTestData[responses.TokenResponse]("ciam_getAuth_res.json")
	ciamGetUserRes := utils.LoadTestData[responses.GraphUserCollection]("ciam_getUser_success_res.json")

	expectedRes := utils.LoadTestData[responses.ApiResponse[any]]("lbe11_withdrawUser_res.json")

	rlpProfileUrl := strings.ReplaceAll(services.ProfileURL, ":api_key", config.GetConfig().Api.Rlp.Core.ApiKey)

	tests := []struct {
		name                 string
		externalId           string
		setupMocks           func()
		expectedHTTPCode     int
		expectedResponseBody any
	}{
		{
			name:       "SUCCESS - User withdrawn",
			externalId: "25052300047",
			setupMocks: func() {
				// Mock RLP Get Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Get(rlpProfileUrl).
					Reply(200).
					JSON(rlpGetProfileRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// CIAM user exists
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", utils.BuildCiamEmailFilter(rlpGetProfileRes.User.Email)).
					Reply(200).
					JSON(ciamGetUserRes)

				// Mock RLP Put Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(rlpProfileUrl).
					Reply(200).
					JSON(rlpUpdateProfileRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM update user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Patch(fmt.Sprintf("%s/%s", services.CiamUserURL, ciamGetUserRes.Value[0].ID)).
					Reply(204)

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
			expectedHTTPCode:     http.StatusOK,
			expectedResponseBody: expectedRes,
		},
		{
			name:       "CONFLICT - RLP Get user not found",
			externalId: "25052300047",
			setupMocks: func() {
				// Mock RLP Get Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Get(rlpProfileUrl).
					Reply(400).
					JSON(rlpUserProfileNotFoundRes)
			},
			expectedHTTPCode:     http.StatusConflict,
			expectedResponseBody: responses.ExistingUserNotFoundErrorResponse(),
		},
		{
			name:       "CONFLICT - CIAM get user by email not found",
			externalId: "25052300047",
			setupMocks: func() {
				// Mock RLP Get Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Get(rlpProfileUrl).
					Reply(200).
					JSON(rlpGetProfileRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// CIAM user DOES NOT exist
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", utils.BuildCiamEmailFilter(rlpGetProfileRes.User.Email)).
					Reply(200).
					JSON(map[string]any{"value": []any{}})
			},
			expectedHTTPCode:     http.StatusConflict,
			expectedResponseBody: responses.ExistingUserNotFoundErrorResponse(),
		},
		{
			name:       "CONFLICT - RLP update user profile not found",
			externalId: "25052300047",
			setupMocks: func() {
				// Mock RLP Get Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Get(rlpProfileUrl).
					Reply(200).
					JSON(rlpGetProfileRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// CIAM user exists
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", utils.BuildCiamEmailFilter(rlpGetProfileRes.User.Email)).
					Reply(200).
					JSON(ciamGetUserRes)

				// Mock RLP Put Profile not found
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(rlpProfileUrl).
					Reply(400).
					JSON(rlpUserProfileNotFoundRes)
			},
			expectedHTTPCode:     http.StatusConflict,
			expectedResponseBody: responses.ExistingUserNotFoundErrorResponse(),
		},
		{
			name:       "ERROR - RLP get user failed",
			externalId: "25052300047",
			setupMocks: func() {
				// Mock RLP Get Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Get(rlpProfileUrl).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:       "ERROR - CIAM get user by email failed",
			externalId: "25052300047",
			setupMocks: func() {
				// Mock RLP Get Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Get(rlpProfileUrl).
					Reply(200).
					JSON(rlpGetProfileRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// CIAM user exists error
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", utils.BuildCiamEmailFilter(rlpGetProfileRes.User.Email)).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:       "ERROR - RLP put profile withdraw failed",
			externalId: "25052300047",
			setupMocks: func() {
				// Mock RLP Get Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Get(rlpProfileUrl).
					Reply(200).
					JSON(rlpGetProfileRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// CIAM user exists
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", utils.BuildCiamEmailFilter(rlpGetProfileRes.User.Email)).
					Reply(200).
					JSON(ciamGetUserRes)

				// Mock RLP Put Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(rlpProfileUrl).
					Reply(200).
					JSON(rlpUpdateProfileRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM update user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Patch(fmt.Sprintf("%s/%s", services.CiamUserURL, ciamGetUserRes.Value[0].ID)).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:       "ERROR - CIAM update user withdraw failed",
			externalId: "25052300047",
			setupMocks: func() {
				// Mock RLP Get Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Get(rlpProfileUrl).
					Reply(200).
					JSON(rlpGetProfileRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// CIAM user exists
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", utils.BuildCiamEmailFilter(rlpGetProfileRes.User.Email)).
					Reply(200).
					JSON(ciamGetUserRes)

				// Mock RLP Put Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(rlpProfileUrl).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:       "ERROR - ACS send email failed",
			externalId: "25052300047",
			setupMocks: func() {
				// Mock RLP Get Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Get(rlpProfileUrl).
					Reply(200).
					JSON(rlpGetProfileRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// CIAM user exists
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", utils.BuildCiamEmailFilter(rlpGetProfileRes.User.Email)).
					Reply(200).
					JSON(ciamGetUserRes)

				// Mock RLP Put Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(rlpProfileUrl).
					Reply(200).
					JSON(rlpUpdateProfileRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM update user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Patch(fmt.Sprintf("%s/%s", services.CiamUserURL, ciamGetUserRes.Value[0].ID)).
					Reply(204)

				// Mock ACS auth
				gock.New(config.GetConfig().Api.Acs.Host).
					Post(services.AcsAuthURL).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock ACS email send
				sendEndpoint := strings.ReplaceAll(
					services.AcsSendEmailByTemplateURL,
					":template_name", services.AcsEmailTemplateRequestOtp, //TODO: update proper impl
				)
				gock.New(config.GetConfig().Api.Acs.Host).
					Post(sendEndpoint).
					Reply(500).
					JSON(nil)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			// Setup mocks
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodPut, "/user/archive/"+tt.externalId, nil)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedHTTPCode, rec.Code)

			if rec.Code == http.StatusOK || rec.Code == http.StatusConflict {
				var resp responses.ApiResponse[any]
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)

				if tt.expectedResponseBody != nil {
					expected, _ := json.Marshal(tt.expectedResponseBody)
					actual, _ := json.Marshal(resp)
					assert.JSONEq(t, string(expected), string(actual))
				}
			}
		})
	}
}
