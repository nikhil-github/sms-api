package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

// Message represent input payload.
type Message struct {
	PhoneNumber string   `json:"phone_number" `
	Texts       []string `json:"texts"`
}

// Result represent status of sms send request.
type Result struct {
	Status []string `json:"status"`
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

// Send handles incoming request to send sms.
// POST /api/v1/sms/send
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

		if err := validate(m); err != nil {
			responseBadRequest(w, enc, err.Error())
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

		var status []string
		for _, text := range m.Texts {
			if len(text) == 0 {
				continue
			}
			err := sender.Send(ctx, number, text)
			if err != nil {
				status = append(status, "failed")
				logger.Error("Unable to send sms", zap.String("text", text), zap.Error(err))
			} else {
				status = append(status, "success")
			}
		}
		responseOK(w, enc, status)
		return
	}
}

func validate(m Message) error {
	if m.PhoneNumber == "" {
		return errors.New("phone number missing")
	}
	if len(m.Texts) == 0 {
		return errors.New("texts missing")
	}
	if len(m.Texts) > 3 {
		return errors.New("max allowed text count is 3")
	}
	for _, t := range m.Texts {
		if len(t) > 160 {
			return errors.New("max allowed text length is 160")
		}
	}
	return nil
}

func responseOK(w http.ResponseWriter, encoder *json.Encoder, response []string) {
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
func NewStatus(status []string) *Result {
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
