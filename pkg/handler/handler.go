package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// Message represent payload.
type Message struct {
	PhoneNumber string   `json:"phone_number" `
	Texts       []string `json:"texts"`
}

// Result represent status of sms send request.
type Result struct {
	Status map[int]string `json:"status"`
}

// ErrorMsg represent error msg.
type ErrorMsg struct {
	Message string `json:"message"`
}

// Sender provides method to send sms.
type Sender interface {
	Send(ctx context.Context, phoneNumber int64, text string) error
}

// Formatter provides method to format phone number.
type Formatter interface {
	Format(ctx context.Context, phoneNumber string) (int64, bool, error)
}

// Send validates the input and call service to send sms.
func Send(logger *zap.Logger, sender Sender, formatter Formatter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		enc := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var m Message
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			logger.Warn("Unable to parse JSON from request body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if msg := validate(m); msg != "" {
			responseBadRequest(w, enc, msg)
			return
		}

		number, valid, err := formatter.Format(ctx, m.PhoneNumber)
		if err != nil {
			logger.Error("Unable to validate number", zap.Error(err))
			serverError(w, enc, "unable to validate number")
			return
		}
		if !valid {
			responseBadRequest(w, enc, "invalid phone number")
			return
		}

		status := make(map[int]string, 3)
		for i, text := range m.Texts {
			err := sender.Send(ctx, number, text)
			if err != nil {
				status[i] = "failed"
				logger.Error("Unable to send sms", zap.Error(err))
			} else {
				status[i] = "success"
			}
		}
		responseOK(w, enc, status)
		return
	}
}

func validate(m Message) string {
	if m.PhoneNumber == "" {
		return "phone number missing"
	}
	if len(m.Texts) == 0 {
		return "texts missing"
	}
	if len(m.Texts) > 3 {
		return "max allowed text count is 3"
	}
	for _, t := range m.Texts {
		if !validateLength(t) {
			return "max allowed text length is 160"
		}
	}
	return ""
}

func validateLength(t string) bool {
	if len(t) > 160 {
		return false
	}
	return true
}

func responseOK(w http.ResponseWriter, encoder *json.Encoder, response map[int]string) {
	w.WriteHeader(http.StatusOK)
	encoder.Encode(NewStatus(response))
}

func responseBadRequest(w http.ResponseWriter, encoder *json.Encoder, response string) {
	w.WriteHeader(http.StatusBadRequest)
	encoder.Encode(NewErrorMsg(response))
}

func serverError(w http.ResponseWriter, encoder *json.Encoder, response string) {
	code := http.StatusInternalServerError
	w.WriteHeader(code)
	encoder.Encode(NewErrorMsg(response))
}

// NewStatus new error message.
func NewStatus(status map[int]string) *Result {
	return &Result{
		Status: status,
	}
}

// NewErrorMsg new error message.
func NewErrorMsg(message string) *ErrorMsg {
	return &ErrorMsg{
		Message: message,
	}
}
