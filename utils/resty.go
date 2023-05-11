package utils

import (
	"github.com/go-resty/resty/v2"
)

var restyClient *resty.Client

func GetRestyClient() *resty.Client {
	if restyClient != nil {
		return restyClient
	}
	client := resty.New()
	client.SetContentLength(true)

	restyClient = client
	return restyClient
}
