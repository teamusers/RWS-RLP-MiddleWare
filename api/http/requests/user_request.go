package requests

type UpdateBurnPinRequest struct {
	BurnPin   *string `json:"burn_PIN" binding:"required,len=4,numeric"`
	DebugMode *int    `json:"debugMode"` //to be removed!
}
