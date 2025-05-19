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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestVerifyUserExistence_Success(t *testing.T) {
	defer gock.Off() // flush mocks after test
	cfg := config.GetConfig()
	sampleEmail := "user@example.com"

	// === 1. Mock CIAM user lookup ===
	accessTokenEndpoint := fmt.Sprintf("/%s%s", cfg.Api.Eeid.TenantID, services.CiamAuthURL)
	gock.New(cfg.Api.Eeid.AuthHost).
		Post(accessTokenEndpoint).
		Reply(200).
		JSON(responses.AcsAuthResponseData{
			AccessToken: "mockToken",
		})

	filter := fmt.Sprintf("mail eq '%s'", sampleEmail)
	gock.New(cfg.Api.Eeid.Host).
		Get(services.CiamUserURL).
		MatchParam("$filter", filter).
		Reply(200).
		JSON(map[string]any{
			"value": []any{}, // empty list = user not found
		})

	// === 2. Mock ACS email sending ===
	gock.New(cfg.Api.Acs.Host).
		Post(services.AcsAuthURL).
		Reply(200).
		JSON(responses.AcsAuthResponseData{
			AccessToken: "mockToken",
		})

	sendEmailEndpoint := strings.ReplaceAll(services.AcsSendEmailByTemplateURL, ":template_name", services.AcsEmailTemplateRequestOtp)
	gock.New(cfg.Api.Acs.Host).
		Post(sendEmailEndpoint).
		Reply(200).
		JSON(nil)

	// === 3. Setup Gin router ===
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/user/register/verify", user.VerifyUserExistence)

	// === 4. Prepare request body ===
	requestBody := requests.VerifyUserExistence{
		Email: sampleEmail,
	}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/user/register/verify", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// === 5. Perform the request ===
	router.ServeHTTP(rec, req)

	// === 6. Assertions ===
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp responses.ApiResponse[model.Otp]
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, codes.SUCCESSFUL, resp.Code)
	assert.Equal(t, "existing user not found", resp.Message)
	assert.NotNil(t, *resp.Data.Otp)
}
