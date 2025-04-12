package responses

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
}
