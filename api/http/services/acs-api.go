package services

import (
	"context"
	"fmt"
	"lbe/api/http/responses"
	"lbe/config"
	"lbe/model"
	"lbe/utils"
	"log"
	"net/http"
	"strings"
)

const (
	//endpoints
	AcsAuthURL                = "/api/v1/auth"
	AcsSendEmailByTemplateURL = "/api/v1/send/template/:template_name"

	// subjects
	AcsEmailSubjectRequestOtp = "RWS Loyalty Program - Verify OTP"

	// template names
	AcsEmailTemplateRequestOtp = "request_email_otp"
)

func getAcsAccessToken(ctx context.Context, client *http.Client) (string, error) {
	appId := config.GetConfig().Api.Acs.AppId
	secretKey := config.GetConfig().Api.Acs.Secret
	reqBody, err := GenerateSignature(appId, secretKey)

	if err != nil {
		log.Printf("unable to generate auth signature: %v", err)
		return "", err
	}

	headers := map[string]string{
		"AppID": config.GetConfig().Api.Acs.AppId,
	}

	if response, err := utils.DoAPIRequest[responses.ApiResponse[responses.AcsAuthResponseData]](model.APIRequestOptions{
		Method:         http.MethodPost,
		URL:            buildFullAcsUrl(AcsAuthURL),
		Body:           reqBody,
		ExpectedStatus: http.StatusOK,
		Headers:        headers,
		Client:         client,
		Context:        ctx,
		ContentType:    model.ContentTypeJson,
	}); err != nil {
		return "", err
	} else {
		return response.Data.AccessToken, nil
	}
}

func PostAcsSendEmailByTemplate(ctx context.Context, client *http.Client, templateName string, payload any) error {
	bearerToken, err := getAcsAccessToken(ctx, client)
	if err != nil {
		log.Printf("error getting acs token: %v", err)
		return err
	}
	headers := map[string]string{
		"AppID": config.GetConfig().Api.Acs.AppId,
	}

	url := strings.ReplaceAll(AcsSendEmailByTemplateURL, ":template_name", templateName)

	if _, err := utils.DoAPIRequest[struct{}](model.APIRequestOptions{
		Method:         http.MethodPost,
		URL:            buildFullAcsUrl(url),
		Body:           payload,
		BearerToken:    bearerToken,
		ExpectedStatus: http.StatusOK,
		Headers:        headers,
		Client:         client,
		Context:        ctx,
		ContentType:    model.ContentTypeJson,
	}); err != nil {
		return err
	}

	return nil
}

func buildFullAcsUrl(endpoint string) string {
	return fmt.Sprintf("%s%s", config.GetConfig().Api.Acs.Host, endpoint)
}
