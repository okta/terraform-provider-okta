package sdk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/peterhellberg/link"
)

// ApiSupplement not all APIs are supported by okta-sdk-golang, this will act as a supplement to the Okta SDK
type ApiSupplement struct {
	BaseURL         string
	Client          *http.Client
	Token           string
	RequestExecutor *okta.RequestExecutor
}

func (m *ApiSupplement) GetSAMLMetdata(id, keyID string) ([]byte, error) {
	return m.GetXml(fmt.Sprintf("%s/api/v1/apps/%s/sso/saml/metadata?kid=%s", m.BaseURL, id, keyID))
}

func (m *ApiSupplement) GetSAMLIdpMetdata(id string) ([]byte, error) {
	return m.GetXml(fmt.Sprintf("%s/api/v1/idps/%s/metadata.xml", m.BaseURL, id))
}

func (m *ApiSupplement) GetXml(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", m.Token))
	req.Header.Add("User-Agent", "Terraform Okta Provider")
	req.Header.Add("Accept", "application/xml")
	resp, err := m.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
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
