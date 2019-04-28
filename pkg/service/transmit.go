package service

const (
	baseURL      = "https://api.transmitsms.com"
	sendSMS      = "/send-sms.json"
	formatNumber = "/format-number.json"
)

// Format represent the format API response.
type Format struct {
	Number Number `json:"number"`
	Error  Error  `json:"error"`
}

// Number represent phone number details.
type Number struct {
	International int64 `json:"international"`
	IsValid       bool  `json:"isValid"`
}

// Response represents send sms response.
type Response struct {
	Error Error `json:"error"`
}

// Error represent transmit API error.
type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
