package utils

import (
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type MicrolinkClient struct {
	HostURL string `json:"url"`
	Resty   *resty.Client
	ApiKey  string `json:"api_token"`
}

func (m *MicrolinkClient) SetHost(host string) {
	m.HostURL = host
}

func (m *MicrolinkClient) SetApiKey(key string) {
	m.ApiKey = key
}

func (m *MicrolinkClient) GetMetaData(url string) *resty.Response {
	res, _ := m.Resty.R().
		SetHeader("x-api-key", m.ApiKey).
		SetResult(&MicrolinkResponse{}).
		SetQueryParam("url", url).
		Get(m.HostURL)
	return res
}

type MicrolinkDataLogo struct {
	URL string `json:"url"`
}

type MicrolinkData struct {
	Title       string            `json:"title"`
	URL         string            `json:"url"`
	Description string            `json:"description"`
	Publisher   string            `json:"publisher"`
	Date        string            `json:"date"`
	Logo        MicrolinkDataLogo `json:"logo"`
}

type MicrolinkResponse struct {
	Status string        `json:"status"`
	Data   MicrolinkData `json:"data"`
}

var microlinkClient *MicrolinkClient

func GetMicrolinkClient() *MicrolinkClient {
	if microlinkClient != nil {
		return microlinkClient
	}
	client := &MicrolinkClient{}
	restyClient := resty.New()
	client.Resty = restyClient
	apiKey := viper.GetString("MicrolinkApikey")
	client.SetApiKey(apiKey)
	client.SetHost("https://pro.microlink.io")
	microlinkClient = client
	return microlinkClient
}
