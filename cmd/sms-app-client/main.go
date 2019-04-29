package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	sendSMSURL = "http://localhost:3001/api/v1/sms/send"
)

func main() {
	ctx := context.Background()
	var netClient = &http.Client{
		Timeout: time.Second * 5,
	}
	sendSMS(ctx, netClient)
}

func sendSMS(ctx context.Context, client *http.Client) {
	r, err := http.NewRequest("POST", sendSMSURL, strings.NewReader(`{"phone_number":"011010101011","texts":["text1"]}`))
	if err != nil {
		log.Fatalf("request creation failed")
	}
	res, err := client.Do(r.WithContext(ctx))
	if err != nil {
		fmt.Println(err)
		log.Fatalf("service call failure %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("failed to send sms http status code %d", res.StatusCode)
	}
	log.Println("send sms successful!!")
}
