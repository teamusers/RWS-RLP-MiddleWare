package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	config "rlp-middleware/config"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	SmtpConfig *config.SmtpConfig
	Dialer     *gomail.Dialer
}

type TemplateData struct {
	Email string
	OTP   string
}

func NewEmailService(smtpConfig *config.SmtpConfig) *EmailService {
	dialer := gomail.NewDialer(
		smtpConfig.Host,
		smtpConfig.Port,
		smtpConfig.User,
		smtpConfig.Password,
	)

	return &EmailService{
		SmtpConfig: smtpConfig,
		Dialer:     dialer,
	}
}

func (es *EmailService) SendOtpEmail(recipient string, templateData TemplateData) error {
	return es.sendEmail(recipient, "RWS Loyalty Program - Verify OTP", "request_email_otp.html", templateData)
}

func (es *EmailService) sendEmail(recipient, subject, templateName string, templateData any) error {
	content, err := es.loadTemplate(templateName, templateData)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	// Create a new message
	m := gomail.NewMessage()

	// Set the sender, recipient, subject, and body of the email
	m.SetHeader("From", es.SmtpConfig.From)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	// Send the email
	err = es.Dialer.DialAndSend(m)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("email sent successfully to %s", recipient)
	return nil
}

func (es *EmailService) loadTemplate(templateName string, data any) (string, error) {
	templatePath := fmt.Sprintf("templates/%s", templateName)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var tplBuffer bytes.Buffer
	err = tmpl.Execute(&tplBuffer, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return tplBuffer.String(), nil
}
