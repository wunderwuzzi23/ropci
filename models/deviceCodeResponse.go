package models

type DeviceCodeResponse struct {
	UserCode        string `json:"user_code"`
	DeviceCode      string `json:"device_code"`
	VerificationUrl string `json:"verification_url"`
	Message         string `json:"message"`
	Interval        int    `json:"interval"`
}
