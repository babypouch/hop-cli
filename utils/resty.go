package utils

import (
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

var restyClient *resty.Client

func GetRestyClient() *resty.Client {
	if restyClient != nil {
		return restyClient
	}
	client := resty.New()
	authToken := viper.GetString("AuthToken")
	client.SetAuthToken(authToken)
	client.SetContentLength(true)

	restyClient = client
	return restyClient
}
