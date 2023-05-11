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
	client.SetAuthToken("06eb9b2c5fb8fd976159f79d3a115f2789a14eaf276414f29cc1e4ad5108ee5cd14c18978330391cad6ddc98360c512327b00bc072e0e9ed3ff4690daefd88cf771326288ec2df4e008df8ab508282d66e8f537045a0703ba25a32ffdf3343ff0eb675d4568b40fe77b3a43767629e0d1b40ff5a94b0f72876d7ec72cfff8f26")
	client.SetContentLength(true)

	restyClient = client
	return restyClient
}
