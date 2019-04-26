package service

const (
	baseURL      = "https://api.transmitsms.com"
	sendSMS      = "/send-sms.json"
	formatNumber = "/format-number.json"
)

type Format struct {
	Number Number `json:"number"`
	Error  Error  `json:"error"`
}

type Number struct {
	International int64 `json:"international"`
	IsValid       bool  `json:"isValid"`
}

type Response struct {
	Error Error `json:"error"`
}

type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
