package utils

import (
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type SerpClient struct {
	HostURL   string `json:"url"`
	Resty     *resty.Client
	AuthToken string `json:"auth_token"`
}

type SerpParams struct {
	Q         string `json:"q"`
	Domain    string `json:"domain"`
	Lang      string `json:"lang"`
	Device    string `json:"device"`
	SerpType  string `json:"serp_type"`
	Loc       string `json:"loc"`
	LocId     string `json:"loc_id"`
	Verbatim  string `json:"verbatim"`
	Gfilter   string `json:"gfilter"`
	Page      string `json:"page"`
	NumResult string `json:"num_result"`
	ApiToken  string `json:"api_token"`
}

func (s *SerpClient) SetHost(url string) {
	s.HostURL = url
}

func (s *SerpClient) SetAuthToken(token string) {
	s.AuthToken = token
}

func (s *SerpClient) GetLive(params map[string]string) *resty.Response {
	res, _ := s.Resty.R().
		SetResult(&SerpLiveResponse{}).
		SetQueryParams(params).
		Get(s.HostURL)

	return res
}

type SerpSearchMetadata struct {
	Id          string
	Status      string
	CreatedAt   string
	ProcessedAt string
}

type SerpSearchResults struct {
	Organic []SerpSearchResultsOrganic `json:"organic"`
}

type SerpSearchResultsOrganic struct {
	Position      string `json:"position"`
	Title         string `json:"title"`
	Link          string `json:"link"`
	DisplayedLink string `json:"displayed_link"`
	CachedPage    string `json:"cached_page"`
	Snippet       string `json:"snippet"`
}

type SerpSearchResponseResults struct {
	SearchMetadata SerpSearchMetadata `json:"search_metadata"`
	Results        SerpSearchResults  `json:"results"`
}

type SerpLiveResponse struct {
	Status  string                    `json:"status"`
	Msg     string                    `json:"msg"`
	Results SerpSearchResponseResults `json:"results"`
}

func (serpLiveRes *SerpLiveResponse) GetOrganic() []SerpSearchResultsOrganic {
	return serpLiveRes.Results.Results.Organic
}

func (serpLiveRes *SerpLiveResponse) FindResultByTitle(title string) *SerpSearchResultsOrganic {
	for _, result := range serpLiveRes.Results.Results.Organic {
		if result.Title == title {
			return &result
		}
	}
	return &SerpSearchResultsOrganic{}
}

var serpClient *SerpClient

func GetSerpClient() *SerpClient {
	if serpClient != nil {
		return serpClient
	}
	client := &SerpClient{}
	hostURL := "https://api.serphouse.com/serp/live"
	authToken := viper.GetString("SerpAuthToken")
	client.SetHost(hostURL)
	client.SetAuthToken(authToken)
	restyClient := resty.New()
	restyClient.SetAuthToken(authToken)
	restyClient.SetContentLength(true)
	client.Resty = restyClient

	serpClient = client
	return serpClient
}
