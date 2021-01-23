package sdk

import (
	"net/url"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/peterhellberg/link"
)

// ApiSupplement not all APIs are supported by okta-sdk-golang, this will act as a supplement to the Okta SDK
type ApiSupplement struct {
	RequestExecutor *okta.RequestExecutor
}

// GetAfterParam grabs after link from link headers if it exists
func GetAfterParam(res *okta.Response) string {
	if res == nil {
		return ""
	}

	linkList := link.ParseHeader(res.Header)
	for _, l := range linkList {
		if l.Rel == "next" {
			parsedURL, err := url.Parse(l.URI)
			if err != nil {
				continue
			}
			q := parsedURL.Query()
			return q.Get("after")
		}
	}

	return ""
}
