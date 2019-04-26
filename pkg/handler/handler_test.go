package handler_test

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/nikhil-github/sms-api/pkg/wiring"
)

func TestSend(t *testing.T) {
	type args struct {
		Input io.Reader
	}
	type fields struct {
		MockExpectations func(m *mockFormatter, s *mockSender)
	}
	type want struct {
		Status int
		Body   string
	}
	testTable := []struct {
		Name   string
		Args   args
		Fields fields
		Want   want
	}{
		{
			Name:   "Failure - missing phone number",
			Args:   args{Input: strings.NewReader(`{"texts":["text"]}`)},
			Fields: fields{MockExpectations: func(m *mockFormatter, s *mockSender) {}},
			Want:   want{Status: http.StatusBadRequest, Body: `{"message": "phone number missing"}`},
		},
		{
			Name:   "Failure - texts count great than 3",
			Args:   args{Input: strings.NewReader(`{"phone_number":"10101010","texts":["text1","text2","text3","text4"]}`)},
			Fields: fields{MockExpectations: func(m *mockFormatter, s *mockSender) {}},
			Want:   want{Status: http.StatusBadRequest, Body: `{"message": "max allowed text count is 3"}`},
		},
		{
			Name: "Failure - Invalid phone number",
			Args: args{Input: strings.NewReader(`{"phone_number":"wrong-number","texts":["text"]}`)},
			Fields: fields{MockExpectations: func(m *mockFormatter, s *mockSender) {
				m.OnFormat("wrong-number").Return(int64(0), false, nil)
			}},
			Want: want{Status: http.StatusBadRequest, Body: `{"message": "invalid phone number"}`},
		},
		{
			Name: "Partial Failure - send sms",
			Args: args{Input: strings.NewReader(`{"phone_number":"wrong-number","texts":["text1","text2"]}`)},
			Fields: fields{MockExpectations: func(m *mockFormatter, s *mockSender) {
				m.OnFormat("wrong-number").Return(int64(88787878), true, nil)
				s.OnSend(int64(88787878), "text1").Return(errors.New("error"))
				s.OnSend(int64(88787878), "text2").Return(nil)
			}},
			Want: want{Status: http.StatusOK, Body: `{"status": {"0": "failed","1": "success"}}`},
		},
		{
			Name: "Success - send sms",
			Args: args{Input: strings.NewReader(`{"phone_number":"wrong-number","texts":["text1"]}`)},
			Fields: fields{MockExpectations: func(m *mockFormatter, s *mockSender) {
				m.OnFormat("wrong-number").Return(int64(5555555), true, nil)
				s.OnSend(int64(5555555), "text1").Return(nil)
			}},
			Want: want{Status: http.StatusOK, Body: `{"status": {"0": "success"}}`},
		},
	}
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			logger := zap.NewNop()
			var m mockFormatter
			var s mockSender

			tt.Fields.MockExpectations(&m, &s)
			params := new(wiring.Params)
			params.Formatter = &m
			params.Sender = &s
			params.Logger = logger
			mx := wiring.NewRouter(params)
			ts := httptest.NewServer(mx)
			defer ts.Close()
			res, err := http.Post(ts.URL+"/api/v1/sms/send", "application/json", tt.Args.Input)
			assert.NoError(t, err, "Error executing request")
			defer res.Body.Close()
			m.AssertExpectations(t)
			assert.Equal(t, tt.Want.Status, res.StatusCode, "status")
			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err, "Error reading response")
			if tt.Want.Body != "" {
				assert.JSONEq(t, tt.Want.Body, string(body), "response")
			}
		})
	}
}

type mockFormatter struct {
	mock.Mock
}

func (m *mockFormatter) Format(ctx context.Context, phoneNumber string) (int64, bool, error) {
	args := m.Called(ctx, phoneNumber)
	return args.Get(0).(int64), args.Get(1).(bool), args.Error(2)
}

func (m *mockFormatter) OnFormat(phoneNumber string) *mock.Call {
	return m.On("Format", mock.AnythingOfType("*context.valueCtx"), phoneNumber)
}

type mockSender struct {
	mock.Mock
}

func (m *mockSender) Send(ctx context.Context, phoneNumber int64, text string) error {
	args := m.Called(ctx, phoneNumber, text)
	return args.Error(0)
}

func (m *mockSender) OnSend(phoneNumber int64, text string) *mock.Call {
	return m.On("Send", mock.AnythingOfType("*context.valueCtx"), phoneNumber, text)
}
