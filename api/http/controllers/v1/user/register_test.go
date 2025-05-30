package user_test

import (
	"bytes"
	"encoding/json"
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
	"log"
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
		// Success cases
		{"SUCCESS - Tier A - 1", "1", "Tier A", false},
		{"SUCCESS - Tier B - 12", "12", "Tier B", false},
		{"SUCCESS - Tier B - 18", "18", "Tier B", false},
		{"SUCCESS - Tier C - 13", "13", "Tier C", false},
		{"SUCCESS - Tier C - 14", "14", "Tier C", false},
		{"SUCCESS - Tier C - 19", "19", "Tier C", false},
		{"SUCCESS - Tier C - 20", "20", "Tier C", false},
		{"SUCCESS - Tier C - 25", "25", "Tier C", false},
		{"SUCCESS - Tier C - 26", "26", "Tier C", false},
		{"SUCCESS - Tier D - 15", "15", "Tier D", false},
		{"SUCCESS - Tier D - 16", "16", "Tier D", false},
		{"SUCCESS - Tier D - 21", "21", "Tier D", false},
		{"SUCCESS - Tier D - 27", "27", "Tier D", false},

		// Error cases
		{"ERROR - Empty input", "", "", true},
		{"ERROR - Non-numeric", "abc", "", true},
		{"ERROR - Negative number", "-2", "", true},
		{"ERROR - Zero", "0", "", true},
		{"ERROR - Unrecognized class level", "99", "", true},
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

// LBE 3 Unit Test
func Test_LBE_3_VerifyUserExistence(t *testing.T) {
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
		expectedResponseBody    any
	}{
		{
			name:        "SUCCESS - User not found",
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
					MatchParam("$filter", utils.BuildCiamEmailFilter(email)).
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
			name:        "CONFLICT - User already exists",
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
					MatchParam("$filter", utils.BuildCiamEmailFilter(email)).
					Reply(200).
					JSON(map[string]any{
						"value": []map[string]any{{"id": "abc123"}},
					})
			},
			expectedHTTPCode:     http.StatusConflict,
			expectedResponseBody: responses.ExistingUserFoundErrorResponse(),
		},
		{
			name:        "ERROR - CIAM get user by email fail",
			requestBody: requests.VerifyUserExistence{Email: "existing@example.com"},
			setupMocks: func(email string) {
				// CIAM auth error
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:        "ERROR - ACS send email fail",
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
					MatchParam("$filter", utils.BuildCiamEmailFilter(email)).
					Reply(200).
					JSON(map[string]any{"value": []any{}})

				// Mock ACS auth error
				gock.New(config.GetConfig().Api.Acs.Host).
					Post(services.AcsAuthURL).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},

		{
			name:                 "ERROR - Invalid request body",
			requestBody:          nil, // raw string invalid body
			setupMocks:           func(email string) {},
			expectedHTTPCode:     http.StatusBadRequest,
			expectedResponseBody: responses.InvalidRequestBodyErrorResponse(),
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
				var resp responses.ApiResponse[any]
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)

				if tt.expectedResponseBody != nil {
					expected, _ := json.Marshal(tt.expectedResponseBody)
					actual, _ := json.Marshal(resp)
					assert.JSONEq(t, string(expected), string(actual))
				} else {
					assert.Equal(t, tt.expectedResponseCode, resp.Code)
					assert.Equal(t, tt.expectedResponseMessage, resp.Message)
				}
			}
		})
	}
}

// LBE 4 Unit Test
func Test_LBE_4_CreateUser(t *testing.T) {
	defer gock.Off()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/register", user.CreateUser)

	rlpCreateProfileRes := utils.LoadTestData[responses.GetUserResponse]("rlp_create_profile_res.json")
	rlpUpdateProfileRes := utils.LoadTestData[responses.GetUserResponse]("rlp_update_profile_res.json")
	rlpUpdateProfileGrRes := utils.LoadTestData[responses.GetUserResponse]("rlp_update_profile_GR_res.json")
	rlpUpdateProfileTmRes := utils.LoadTestData[responses.GetUserResponse]("rlp_update_profile_TM_res.json")
	rlpUserProfileNotFoundRes := utils.LoadTestData[responses.UserProfileErrorResponse]("rlp_user_not_found_err_res.json")
	rlpUpdateUserTierEventRes := utils.LoadTestData[any]("rlp_updateUserTier_event_res.json")

	ciamGetAuth := utils.LoadTestData[responses.TokenResponse]("ciam_getAuth_res.json")
	ciamRegisterUserRes := utils.LoadTestData[responses.GraphCreateUserResponse]("ciam_createUser_success_res.json")
	ciamRegisterUserErrorRes := utils.LoadTestData[responses.GraphApiErrorResponse]("ciam_createUser_error_res.json")

	validSampleReqNew := utils.LoadTestData[requests.RegisterUser]("lbe4_createUser_NEW_req.json")
	validSampleReqGrCms := utils.LoadTestData[requests.RegisterUser]("lbe4_createUser_GRCMS_req.json")
	validSampleReqGr := utils.LoadTestData[requests.RegisterUser]("lbe4_createUser_GR_req.json")
	invalidSampleReqGr := utils.LoadTestData[requests.RegisterUser]("lbe4_createUser_invalidGrClass_req.json")
	validSampleReqTm := utils.LoadTestData[requests.RegisterUser]("lbe4_createUser_TM_req.json")

	expectedResNew := utils.LoadTestData[responses.ApiResponse[any]]("lbe4_createUser_NEW_res.json")
	expectedResGr := utils.LoadTestData[responses.ApiResponse[any]]("lbe4_createUser_GR_res.json")
	expectedResTm := utils.LoadTestData[responses.ApiResponse[any]]("lbe4_createUser_TM_res.json")

	createRlpUserProfileUrl := strings.ReplaceAll(services.CreateProfileURL, ":api_key", config.GetConfig().Api.Rlp.Core.ApiKey)
	rlpProfileUrl := strings.ReplaceAll(services.ProfileURL, ":api_key", config.GetConfig().Api.Rlp.Core.ApiKey)
	updateRlpUserProfileUrl := fmt.Sprintf("%s/.+", rlpProfileUrl)

	tests := []struct {
		name                 string
		requestBody          any
		setupMocks           func(grId, email string)
		expectedHTTPCode     int
		expectedResponseBody any
	}{
		{
			name:        "SUCCESS - NEW user registration",
			requestBody: validSampleReqNew,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM create user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Post(services.CiamUserURL).
					Reply(201).
					JSON(ciamRegisterUserRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM add schema extensions
				gock.New(config.GetConfig().Api.Eeid.Host).
					Patch(fmt.Sprintf("%s/%s", services.CiamUserURL, ciamRegisterUserRes.Id)).
					Reply(204)

				// Mock RLP Initial User Creation
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Post(createRlpUserProfileUrl).
					Reply(200).
					JSON(rlpCreateProfileRes)

				// Mock RLP Put Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(updateRlpUserProfileUrl).
					Reply(200).
					JSON(rlpUpdateProfileRes)

				// Mock RLP Update User Tier Event
				gock.New(config.GetConfig().Api.Rlp.Offers.Host).
					Post(services.EventUrl).
					Reply(200).
					JSON(rlpUpdateUserTierEventRes)
			},
			expectedHTTPCode:     http.StatusCreated,
			expectedResponseBody: expectedResNew,
		},
		{
			name:        "SUCCESS - GR CMS user registration",
			requestBody: validSampleReqGrCms,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM create user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Post(services.CiamUserURL).
					Reply(201).
					JSON(ciamRegisterUserRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM add schema extensions
				gock.New(config.GetConfig().Api.Eeid.Host).
					Patch(fmt.Sprintf("%s/%s", services.CiamUserURL, ciamRegisterUserRes.Id)).
					Reply(204)

				// Mock RLP Initial User Creation
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Post(createRlpUserProfileUrl).
					Reply(200).
					JSON(rlpCreateProfileRes)

				// Mock RLP Put Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(updateRlpUserProfileUrl).
					Reply(200).
					JSON(rlpUpdateProfileGrRes)

				// Mock RLP Update User Tier Event
				gock.New(config.GetConfig().Api.Rlp.Offers.Host).
					Post(services.EventUrl).
					Reply(200).
					JSON(rlpUpdateUserTierEventRes)
			},
			expectedHTTPCode:     http.StatusCreated,
			expectedResponseBody: expectedResGr,
		},
		{
			name:        "SUCCESS - GR user registration",
			requestBody: validSampleReqGr,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM create user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Post(services.CiamUserURL).
					Reply(201).
					JSON(ciamRegisterUserRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM add schema extensions
				gock.New(config.GetConfig().Api.Eeid.Host).
					Patch(fmt.Sprintf("%s/%s", services.CiamUserURL, ciamRegisterUserRes.Id)).
					Reply(204)

				// Mock RLP Initial User Creation
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Post(createRlpUserProfileUrl).
					Reply(200).
					JSON(rlpCreateProfileRes)

				// Mock RLP Put Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(updateRlpUserProfileUrl).
					Reply(200).
					JSON(rlpUpdateProfileGrRes)

				// Mock RLP Update User Tier Event
				gock.New(config.GetConfig().Api.Rlp.Offers.Host).
					Post(services.EventUrl).
					Reply(200).
					JSON(rlpUpdateUserTierEventRes)
			},
			expectedHTTPCode:     http.StatusCreated,
			expectedResponseBody: expectedResGr,
		},
		{
			name:        "SUCCESS - TM user registration",
			requestBody: validSampleReqTm,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM create user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Post(services.CiamUserURL).
					Reply(201).
					JSON(ciamRegisterUserRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM add schema extensions
				gock.New(config.GetConfig().Api.Eeid.Host).
					Patch(fmt.Sprintf("%s/%s", services.CiamUserURL, ciamRegisterUserRes.Id)).
					Reply(204)

				// Mock RLP Initial User Creation
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Post(createRlpUserProfileUrl).
					Reply(200).
					JSON(rlpCreateProfileRes)

				// Mock RLP Put Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(updateRlpUserProfileUrl).
					Reply(200).
					JSON(rlpUpdateProfileTmRes)

				// Mock RLP Update User Tier Event
				gock.New(config.GetConfig().Api.Rlp.Offers.Host).
					Post(services.EventUrl).
					Reply(200).
					JSON(rlpUpdateUserTierEventRes)
			},
			expectedHTTPCode:     http.StatusCreated,
			expectedResponseBody: expectedResTm,
		},
		{
			name: "CONFLICT - GR CMS cache not found",
			requestBody: requests.RegisterUser{
				SignUpType: codes.SignUpTypeGRCMS,
				RegId:      "0000",
			},
			setupMocks:           func(grId, email string) {},
			expectedHTTPCode:     http.StatusConflict,
			expectedResponseBody: responses.CachedProfileNotFoundErrorResponse(),
		},
		{
			name:        "CONFLICT - Invalid GR Class",
			requestBody: invalidSampleReqGr,
			setupMocks: func(grId, email string) {
			},
			expectedHTTPCode:     http.StatusConflict,
			expectedResponseBody: responses.InvalidGrMemberClassErrorResponse(),
		},
		{
			name:        "CONFLICT - CIAM user already exists",
			requestBody: validSampleReqNew,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM create user error due to existing user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Post(services.CiamUserURL).
					Reply(400).
					JSON(ciamRegisterUserErrorRes)
			},
			expectedHTTPCode:     http.StatusConflict,
			expectedResponseBody: responses.ExistingUserFoundErrorResponse(),
		},
		{
			name:        "CONFLICT - RLP Update Profile user not found",
			requestBody: validSampleReqNew,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM create user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Post(services.CiamUserURL).
					Reply(201).
					JSON(ciamRegisterUserRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM add schema extensions
				gock.New(config.GetConfig().Api.Eeid.Host).
					Patch(fmt.Sprintf("%s/%s", services.CiamUserURL, ciamRegisterUserRes.Id)).
					Reply(204)

				// Mock RLP Initial User Creation
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Post(createRlpUserProfileUrl).
					Reply(200).
					JSON(rlpCreateProfileRes)

				// Mock RLP Put Profile user does not exist
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(updateRlpUserProfileUrl).
					Reply(400).
					JSON(rlpUserProfileNotFoundRes)
			},
			expectedHTTPCode:     http.StatusConflict,
			expectedResponseBody: responses.ExistingUserNotFoundErrorResponse(),
		},
		{
			name:        "ERROR - CIAM register user fail",
			requestBody: validSampleReqNew,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM create user error
				gock.New(config.GetConfig().Api.Eeid.Host).
					Post(services.CiamUserURL).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:        "ERROR - CIAM patch user schema extensions fail",
			requestBody: validSampleReqNew,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM create user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Post(services.CiamUserURL).
					Reply(201).
					JSON(ciamRegisterUserRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM add schema extensions error
				gock.New(config.GetConfig().Api.Eeid.Host).
					Patch(fmt.Sprintf("%s/%s", services.CiamUserURL, ciamRegisterUserRes.Id)).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:        "ERROR - RLP Create initial user profile fail",
			requestBody: validSampleReqNew,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM create user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Post(services.CiamUserURL).
					Reply(201).
					JSON(ciamRegisterUserRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM add schema extensions
				gock.New(config.GetConfig().Api.Eeid.Host).
					Patch(fmt.Sprintf("%s/%s", services.CiamUserURL, ciamRegisterUserRes.Id)).
					Reply(204)

				// Mock RLP Initial User Creation
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Post(createRlpUserProfileUrl).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:        "ERROR - RLP Update profile fail",
			requestBody: validSampleReqNew,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM create user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Post(services.CiamUserURL).
					Reply(201).
					JSON(ciamRegisterUserRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM add schema extensions
				gock.New(config.GetConfig().Api.Eeid.Host).
					Patch(fmt.Sprintf("%s/%s", services.CiamUserURL, ciamRegisterUserRes.Id)).
					Reply(204)

				// Mock RLP Initial User Creation
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Post(createRlpUserProfileUrl).
					Reply(200).
					JSON(rlpCreateProfileRes)

				// Mock RLP Put Profile fail
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(updateRlpUserProfileUrl).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{ //TODO: add rollback check
			name:        "ERROR - RLP update user tier fail",
			requestBody: validSampleReqNew,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM create user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Post(services.CiamUserURL).
					Reply(201).
					JSON(ciamRegisterUserRes)

				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(ciamGetAuth)

				// Mock CIAM add schema extensions
				gock.New(config.GetConfig().Api.Eeid.Host).
					Patch(fmt.Sprintf("%s/%s", services.CiamUserURL, ciamRegisterUserRes.Id)).
					Reply(204)

				// Mock RLP Initial User Creation
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Post(createRlpUserProfileUrl).
					Reply(200).
					JSON(rlpCreateProfileRes)

				// Mock RLP Put Profile
				gock.New(config.GetConfig().Api.Rlp.Core.Host).
					Put(updateRlpUserProfileUrl).
					Reply(200).
					JSON(rlpUpdateProfileRes)

				// Mock RLP Update User Tier Event Fail
				gock.New(config.GetConfig().Api.Rlp.Offers.Host).
					Post(services.EventUrl).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:                 "ERROR - Invalid request body",
			requestBody:          `{}`,
			setupMocks:           func(email, grId string) {}, // No mock needed
			expectedHTTPCode:     http.StatusBadRequest,
			expectedResponseBody: responses.InvalidRequestBodySpecificErrorResponse("invalid sign_up_type provided"),
		},
		{
			name:                 "ERROR - Invalid JSON ShouldBindJSON",
			requestBody:          `{"user": "invalid-json}`,   // malformed JSON (missing closing quote)
			setupMocks:           func(email, grId string) {}, // No mocks needed
			expectedHTTPCode:     http.StatusBadRequest,
			expectedResponseBody: responses.InvalidRequestBodyErrorResponse(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			//setup cache
			system.ObjectSet(validSampleReqGrCms.RegId, validSampleReqGr.User, 30*time.Minute)
			defer system.ObjectDelete(validSampleReqGrCms.RegId)

			var grId, email string
			if req, ok := tt.requestBody.(requests.RegisterUser); ok {
				if req.SignUpType == codes.SignUpTypeGR {
					if req.User.GrProfile != nil {
						grId = req.User.GrProfile.Id
					}
				}
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

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedHTTPCode, rec.Code)

			if rec.Code == http.StatusCreated || rec.Code == http.StatusConflict {
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

// LBE 6 Unit Test
func Test_LBE_6_VerifyGrExistence(t *testing.T) {
	defer gock.Off()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/user/gr", user.VerifyGrExistence)

	validSampleReq := utils.LoadTestData[requests.VerifyGrUser]("lbe6_verifyGrExistence_req.json")

	tests := []struct {
		name                    string
		requestBody             any
		setupMocks              func(grId, email string)
		expectedHTTPCode        int
		expectedResponseCode    int64
		expectedResponseMessage string
		expectedResponseBody    any
	}{
		{
			name:        "SUCCESS - GR ID not found",
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
					MatchParam("$filter", utils.BuildCiamGrIdFilter(grId)).
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
			expectedResponseMessage: "gr profile found", //TODO: use correct response - need to fix otp impl to be mockable
		},
		{
			name:        "CONFLICT - GR ID already linked",
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
					MatchParam("$filter", utils.BuildCiamGrIdFilter(grId)).
					Reply(200).
					JSON(responses.GraphUserCollection{
						Value: []responses.GraphUser{{ID: "existing-user"}},
					})
			},
			expectedHTTPCode:     http.StatusConflict,
			expectedResponseBody: responses.GrMemberIdLinkedErrorResponse(),
		},
		{
			name:        "ERROR - CIAM get user by grId fail",
			requestBody: validSampleReq,
			setupMocks: func(grId, email string) {
				// Mock CIAM auth
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:        "ERROR - CMS profile fetch",
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
					MatchParam("$filter", utils.BuildCiamGrIdFilter(grId)).
					Reply(200).
					JSON(responses.GraphUserCollection{})

				// CMS member fetch error
				gock.New(config.GetConfig().Api.Cms.Host).
					Get(services.GetMemberURL).
					MatchParam("systemId", config.GetConfig().Api.Cms.SystemID).
					MatchParam("memberId", grId).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:        "ERROR - ACS send email fail",
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
					MatchParam("$filter", utils.BuildCiamGrIdFilter(grId)).
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
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:                 "ERROR - Invalid request body",
			requestBody:          `{}`,
			setupMocks:           func(email, grId string) {}, // No mock needed
			expectedHTTPCode:     http.StatusBadRequest,
			expectedResponseBody: responses.InvalidRequestBodyErrorResponse(),
		},
		{
			name:                 "ERROR - Invalid JSON ShouldBindJSON",
			requestBody:          `{"user": "invalid-json}`,   // malformed JSON (missing closing quote)
			setupMocks:           func(email, grId string) {}, // No mocks needed
			expectedHTTPCode:     http.StatusBadRequest,
			expectedResponseBody: responses.InvalidRequestBodyErrorResponse(),
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

				if tt.expectedResponseBody != nil {
					expected, _ := json.Marshal(tt.expectedResponseBody)
					actual, _ := json.Marshal(resp)
					assert.JSONEq(t, string(expected), string(actual))
				} else {
					assert.Equal(t, tt.expectedResponseCode, resp.Code)
					assert.Equal(t, tt.expectedResponseMessage, resp.Message)
				}
			}
		})
	}
}

// LBE 7 Unit Test
func Test_LBE_7_VerifyGrCmsExistence(t *testing.T) {
	defer gock.Off()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/user/gr-cms", user.VerifyGrCmsExistence)

	validSampleReq := utils.LoadTestData[requests.VerifyGrCmsUser]("lbe7_verifyGrCmsExistence_req.json")

	tests := []struct {
		name                    string
		requestBody             any
		setupMocks              func(email, grId string)
		expectedHTTPCode        int
		expectedResponseCode    int64
		expectedResponseMessage string
		expectedResponseBody    any
	}{
		{
			name:        "SUCCESS - User and GR ID not found",
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
					MatchParam("$filter", utils.BuildCiamEmailFilter(email)).
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
					MatchParam("$filter", utils.BuildCiamGrIdFilter(grId)).
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
			name:        "CONFLICT - Existing email found",
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
					MatchParam("$filter", utils.BuildCiamEmailFilter(email)).
					Reply(200).
					JSON(responses.GraphUserCollection{
						Value: []responses.GraphUser{
							{ID: "abc123"},
						},
					})
			},
			expectedHTTPCode:     http.StatusConflict,
			expectedResponseBody: responses.ExistingUserFoundErrorResponse(),
		},
		{
			name:        "CONFLICT - Existing GR ID found",
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
					MatchParam("$filter", utils.BuildCiamEmailFilter(email)).
					Reply(200).
					JSON(responses.GraphUserCollection{})

				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user by GR ID returns a found user
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", utils.BuildCiamGrIdFilter(grId)).
					Reply(200).
					JSON(responses.GraphUserCollection{
						Value: []responses.GraphUser{
							{ID: "abc123"},
						},
					})
			},
			expectedHTTPCode:     http.StatusConflict,
			expectedResponseBody: responses.GrMemberIdLinkedErrorResponse(),
		},
		{
			name:        "ERROR - CIAM get user by email fail",
			requestBody: validSampleReq,
			setupMocks: func(email, grId string) {

				// Mock CIAM auth error
				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:        "ERROR - CIAM get user by grId fail",
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
					MatchParam("$filter", utils.BuildCiamEmailFilter(email)).
					Reply(200).
					JSON(responses.GraphUserCollection{})

				gock.New(config.GetConfig().Api.Eeid.AuthHost).
					Post(fmt.Sprintf("/%s%s", config.GetConfig().Api.Eeid.TenantID, services.CiamAuthURL)).
					Reply(200).
					JSON(responses.AcsAuthResponseData{AccessToken: "mockToken"})

				// Mock CIAM user by GR ID error
				gock.New(config.GetConfig().Api.Eeid.Host).
					Get(services.CiamUserURL).
					MatchParam("$filter", utils.BuildCiamGrIdFilter(grId)).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:        "ERROR - ACS send email fail",
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
					MatchParam("$filter", utils.BuildCiamEmailFilter(email)).
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
					MatchParam("$filter", utils.BuildCiamGrIdFilter(grId)).
					Reply(200).
					JSON(responses.GraphUserCollection{})

					// Mock ACS auth error
				gock.New(config.GetConfig().Api.Acs.Host).
					Post(services.AcsAuthURL).
					Reply(500)
			},
			expectedHTTPCode:     http.StatusInternalServerError,
			expectedResponseBody: responses.InternalErrorResponse(),
		},
		{
			name:                 "ERROR - Invalid request body",
			requestBody:          nil, // raw string invalid body
			setupMocks:           func(email, grId string) {},
			expectedHTTPCode:     http.StatusBadRequest,
			expectedResponseBody: responses.InvalidRequestBodyErrorResponse(),
		},
		{
			name:                 "ERROR - Invalid JSON ShouldBindJSON",
			requestBody:          `{"user": "invalid-json}`,   // malformed JSON (missing closing quote)
			setupMocks:           func(email, grId string) {}, // No mocks needed
			expectedHTTPCode:     http.StatusBadRequest,
			expectedResponseBody: responses.InvalidRequestBodyErrorResponse(),
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

				if tt.expectedResponseBody != nil {
					expected, _ := json.Marshal(tt.expectedResponseBody)
					actual, _ := json.Marshal(resp)
					assert.JSONEq(t, string(expected), string(actual))
				} else {
					assert.Equal(t, tt.expectedResponseCode, resp.Code)
					assert.Equal(t, tt.expectedResponseMessage, resp.Message)
				}
			}
		})
	}
}

func Test_LBE_8_GetCachedGrCmsProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	expectedRes := utils.LoadTestData[responses.ApiResponse[any]]("lbe8_getCachedGrCmsProfile_res.json")
	var respData responses.VerifyGrCmsUserResponseData

	// Re-marshal the `Data` field to JSON
	dataBytes, err := json.Marshal(expectedRes.Data)
	if err != nil {
		log.Fatalf("failed to marshal Data: %v", err)
	}

	// Unmarshal into the desired struct
	err = json.Unmarshal(dataBytes, &respData)
	if err != nil {
		log.Fatalf("failed to unmarshal into VerifyGrCmsUserResponseData: %v", err)
	}

	cachedRegId := respData.RegId
	cachedUser := &model.User{DateOfBirth: &respData.DateOfBirth}

	tests := []struct {
		name                    string
		regId                   string
		expectedHTTPCode        int
		expectedResponseCode    int64
		expectedResponseMessage string
		expectedResponseBody    any
	}{
		{
			name:                 "SUCCESS - Cache profile found",
			regId:                cachedRegId,
			expectedHTTPCode:     http.StatusOK,
			expectedResponseBody: expectedRes,
		},
		{
			name:                 "CONFLICT - Cache profile not found",
			regId:                "321",
			expectedHTTPCode:     http.StatusConflict,
			expectedResponseBody: responses.CachedProfileNotFoundErrorResponse(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//setup cache
			system.ObjectSet(cachedRegId, cachedUser, 30*time.Minute)
			defer system.ObjectDelete(cachedRegId)

			router := gin.New()
			router.GET("/gr-cms/:reg_id", user.GetCachedGrCmsProfile)

			req := httptest.NewRequest(http.MethodGet, "/gr-cms/"+tt.regId, nil)
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
				} else {
					assert.Equal(t, tt.expectedResponseCode, resp.Code)
					assert.Equal(t, tt.expectedResponseMessage, resp.Message)
				}
			}
		})
	}
}
