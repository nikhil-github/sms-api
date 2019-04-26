package service_test

import (
	"context"
	"net/http"
	"testing"

	"go.uber.org/zap"

	"github.com/nikhil-github/sms-api/pkg/service"
)

func TestFormat(t *testing.T) {
	//ctx := context.Background()
	//authRequest, err := http.NewRequest("POST", "https://oauth.brightcove.com/v4/access_token", bytes.NewReader([]byte("grant_type=client_credentials")))
	//require.NoError(t, err, "Unable to create auth request")
	//authRequest.Header.Add("Content-type", "application/x-www-form-urlencoded")
	//authRequest.Header.Add("Authorization", "Basic Y2xpZW50SUQ6Y2xpZW50S2V5")
	//
	//authRequest1, err := http.NewRequest("POST", "https://oauth.brightcove.com/v4/access_token", bytes.NewReader([]byte("grant_type=client_credentials")))
	//require.NoError(t, err, "Unable to create auth request1")
	//authRequest1.Header.Add("Content-type", "application/x-www-form-urlencoded")
	//authRequest1.Header.Add("Authorization", "Basic Y2xpZW50SUQ6Y2xpZW50S2V5LXdyb25n")

	type args struct {
		ID     string
		Secret string
	}
	type fields struct {
		// MockExpectations func(m *httpClient)
	}
	type want struct {
		Err    string
		Result string
	}
	testTable := []struct {
		Name   string
		Args   args
		Fields fields
		Want   want
	}{
		{
			Name: "Success",
			Want: want{Result: "test-token"},
		},
	}
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			s := service.New("76eaa9a89c1298566e1ce03afd4feec7", "test", http.DefaultClient, zap.NewNop(), "e82183e57166afb3688d4dd962b98937381d396b")
			s.Send(context.Background(), 405990558, "hello http://www.google.com")
			//a, b, c := s.Format(context.Background(), "405990558")
			//fmt.Println(a)
			//fmt.Println(b)
			//fmt.Println(c)
		})
	}
}
