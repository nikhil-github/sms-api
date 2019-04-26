package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/zpnk/go-bitly"
	"go.uber.org/zap"
	"mvdan.cc/xurls"
)

const (
	baseURL      = "https://api.transmitsms.com"
	sendSMS      = "/send-sms.json"
	formatNumber = "/format-number.json"
)

type SenderService struct {
	apiKey      string
	bitlyClient *bitly.Client
	httpClient  HTTPClient
	logger      *zap.Logger
	secret      string
}

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

// HTTPClient an interface for HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// New creates a new SenderService.
func New(apiKey string, secret string, httpClient HTTPClient, l *zap.Logger, bitlyToken string) *SenderService {
	return &SenderService{apiKey: apiKey, secret: secret, httpClient: httpClient, logger: l, bitlyClient: bitly.New(bitlyToken)}
}

// Format method validates and format the given phone number
func (s *SenderService) Format(ctx context.Context, phoneNumber string) (int64, bool, error) {

	data := url.Values{}
	data.Set("msisdn", phoneNumber)
	data.Set("countrycode", "AU")
	req, err := s.request("POST", formatNumber, data.Encode())
	if err != nil {
		return 0, false, err
	}
	res, err := s.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return 0, false, errors.Wrap(err, "failed to format number")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		s.logger.Error("format number error", zap.Int("status code", res.StatusCode))
		return 0, false, fmt.Errorf("error")
	}

	var response Format
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		s.logger.Error("error", zap.Error(err))
		fmt.Println(err)
		return 0, false, fmt.Errorf("error decoding")
	}
	if response.Number.IsValid {
		return response.Number.International, true, nil
	}
	s.logger.Error("error code", zap.String("code", response.Error.Code))
	s.logger.Error("error description", zap.String("description", response.Error.Description))
	return 0, false, err

}

// Send sends sms message to the given number
// links are searched and replaced with short bitly links
func (s *SenderService) Send(ctx context.Context, phoneNumber int64, text string) error {

	text, err := s.replaceLinks(text)
	if err != nil {
		return err
	}
	data := url.Values{}
	data.Set("message", text)
	data.Set("to", strconv.FormatInt(phoneNumber, 10))
	req, err := s.request("POST", sendSMS, data.Encode())
	if err != nil {
		return err
	}
	res, err := s.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return errors.Wrap(err, "failed to send sms")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return nil
	}

	var response Response
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		s.logger.Error("error", zap.Error(err))
		return err
	}

	s.logger.Error("error code", zap.String("code", response.Error.Code))
	s.logger.Error("error description", zap.String("description", response.Error.Description))

	return nil
}

// replaceLinks find links in text and replace them with bitly links
// mvdan.cc/xurls find all links in a string
// github.com/zpnk/go-bitly go bitly client
func (s *SenderService) replaceLinks(text string) (string, error) {
	links := xurls.Strict().FindAllString(text, -1)
	for _, link := range links {
		short, err := s.bitlyClient.Links.Shorten(link)
		if err != nil {
			s.logger.Error("shorten-link-error", zap.Error(err))
			return "", err
		}
		text = strings.Replace(text, link, short.URL, -1)
	}
	fmt.Println("replaced links", text)
	return text, nil
}

func (s *SenderService) request(method string, resource string, data string) (*http.Request, error) {
	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse uri")
	}
	u.Path = resource
	req, err := http.NewRequest(method, u.String(), strings.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req.SetBasicAuth(s.apiKey, s.secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data)))
	return req, nil
}
