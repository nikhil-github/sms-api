package service_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zpnk/go-bitly"
	"go.uber.org/zap"

	"github.com/nikhil-github/sms-api/pkg/service"
)

func TestFormat(t *testing.T) {
	type args struct {
		Number string
	}
	type fields struct {
		MockOperations func(m *httpClient)
	}
	type want struct {
		Err             string
		FormattedNumber int64
		Valid           bool
	}
	testTable := []struct {
		Name   string
		Args   args
		Fields fields
		Want   want
	}{
		{
			Name: "Success : right format",
			Args: args{Number: "1234567890"},
			Fields: fields{MockOperations: func(c *httpClient) {
				data := url.Values{}
				data.Set("msisdn", "1234567890")
				data.Set("countrycode", "AU")
				req, err := http.NewRequest("POST", "https://api.transmitsms.com/format-number.json", strings.NewReader(data.Encode()))
				if err != nil {
					panic(err)
				}
				c.OnDo(req).Return(mockResponse(http.StatusOK, []byte(`{"number":{"international":611234567890,"isValid":true},"error" : {"code":"SUCCESS","description":" "}}`)), nil).Once()
			}},
			Want: want{FormattedNumber: int64(611234567890), Valid: true},
		},
		{
			Name: "Failure : wrong format",
			Args: args{Number: "12345678901"},
			Fields: fields{MockOperations: func(c *httpClient) {
				data := url.Values{}
				data.Set("msisdn", "12345678901")
				data.Set("countrycode", "AU")
				req, err := http.NewRequest("POST", "https://api.transmitsms.com/format-number.json", strings.NewReader(data.Encode()))
				if err != nil {
					panic(err)
				}
				c.OnDo(req).Return(mockResponse(http.StatusOK, []byte(`{"number":{"international":6112345678901,"isValid":false},"error" : {"code":"SUCCESS","description":"OK"}}`)), nil).Once()
			}},
			Want: want{FormattedNumber: int64(0), Valid: false},
		},
		{
			Name: "Error",
			Args: args{Number: "12345678901"},
			Fields: fields{MockOperations: func(c *httpClient) {
				data := url.Values{}
				data.Set("msisdn", "12345678901")
				data.Set("countrycode", "AU")
				req, err := http.NewRequest("POST", "https://api.transmitsms.com/format-number.json", strings.NewReader(data.Encode()))
				if err != nil {
					panic(err)
				}
				c.OnDo(req).Return(mockResponse(http.StatusBadRequest, []byte(`{}`)), errors.New("failed")).Once()
			}},
			Want: want{Err: "failed to format number: failed"},
		},
	}
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			var client httpClient
			tt.Fields.MockOperations(&client)
			s := service.New("key", "secret", &client, zap.NewNop(), bitly.New("dummy"))
			result, valid, err := s.Format(context.Background(), tt.Args.Number)
			client.AssertExpectations(t)
			if tt.Want.Err != "" {
				assert.EqualError(t, err, tt.Want.Err, "error message")
				return
			}
			require.NoError(t, err, "error")
			assert.Equal(t, tt.Want.FormattedNumber, result, "result")
			assert.Equal(t, tt.Want.Valid, valid, "valid")

		})
	}
}

func TestSend(t *testing.T) {
	type args struct {
		Number int64
		Text   string
	}
	type fields struct {
		MockOperations func(m *httpClient, b *mockBitly)
	}
	type want struct {
		Err string
	}
	testTable := []struct {
		Name   string
		Args   args
		Fields fields
		Want   want
	}{
		{
			Name: "Success : send sms",
			Args: args{Number: int64(1234567890), Text: "text"},
			Fields: fields{MockOperations: func(c *httpClient, b *mockBitly) {
				data := url.Values{}
				data.Set("message", "text")
				data.Set("to", strconv.FormatInt(int64(1234567890), 10))
				req, err := http.NewRequest("POST", "https://api.transmitsms.com/send-sms.json", strings.NewReader(data.Encode()))
				if err != nil {
					panic(err)
				}
				c.OnDo(req).Return(mockResponse(http.StatusOK, []byte(`{"error":{"code":"SUCCESS","description":"OK"}}`)), nil).Once()
				b.On("Shorten").Return(bitly.Link{URL: "http://bit.ly/xyz"})
			}},
		},
		{
			Name: "Error : send sms",
			Args: args{Number: int64(1234567890), Text: "text-wrong"},
			Fields: fields{MockOperations: func(c *httpClient, b *mockBitly) {
				data := url.Values{}
				data.Set("message", "text-wrong")
				data.Set("to", strconv.FormatInt(int64(1234567890), 10))
				req, err := http.NewRequest("POST", "https://api.transmitsms.com/send-sms.json", strings.NewReader(data.Encode()))
				if err != nil {
					panic(err)
				}
				c.OnDo(req).Return(mockResponse(http.StatusOK, []byte(`{"error":{"code":"FIELD_INVALID","description":"is not a valid number."}}`)), errors.New("fail")).Once()
				b.On("Shorten").Return(bitly.Link{URL: "http://bit.ly/xyz"})
			}},
			Want: want{Err: "failed to send sms: fail"},
		}}
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			var client httpClient
			var m mockBitly
			tt.Fields.MockOperations(&client, &m)
			s := service.New("key", "secret", &client, zap.NewNop(), bitly.New("dummy"))
			err := s.Send(context.Background(), tt.Args.Number, tt.Args.Text)
			client.AssertExpectations(t)
			if tt.Want.Err != "" {
				assert.EqualError(t, err, tt.Want.Err, "error message")
				return
			}
			require.NoError(t, err, "error")
		})
	}
}

type httpClient struct {
	mock.Mock
}

func (m *httpClient) Do(req *http.Request) (response *http.Response, error error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func mockResponse(statusCode int, body []byte) *http.Response {
	response := httptest.NewRecorder()
	response.WriteHeader(statusCode)
	response.Write(body)
	return response.Result()
}

func (m *httpClient) OnDo(expected *http.Request) *mock.Call {
	return m.On("Do", mock.MatchedBy(func(request *http.Request) bool {
		if expected.Method != request.Method ||
			expected.URL.String() != request.URL.String() {
			return false
		}
		if !reflect.DeepEqual(expected.Body, request.Body) {
			return false
		}
		return true
	}))
}

type mockBitly struct {
	mock.Mock
	*bitly.Client
}

func (m *mockBitly) Shorten(longURL string) (bitly.Link, error) {
	args := m.Called(longURL)
	return args.Get(0).(bitly.Link), args.Error(1)
}
