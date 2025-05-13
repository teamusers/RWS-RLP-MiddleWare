package requests

type AcsSendEmailByTemplateRequest struct {
	Email   string `json:"email,omitempty"`
	Subject string `json:"subject,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type RequestEmailOtpTemplateData struct {
	Email string `json:"email,omitempty"`
	Otp   string `json:"otp,omitempty"`
}
