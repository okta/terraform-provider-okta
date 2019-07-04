package okta

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/okta/okta-sdk-golang/okta"
)

// ApiSupplement not all APIs are supported by okta-sdk-golang, this will act as a supplement to the Okta SDK
type ApiSupplement struct {
	baseURL         string
	client          *http.Client
	token           string
	requestExecutor *okta.RequestExecutor
}

func (m *ApiSupplement) GetSAMLMetdata(id, keyID string) ([]byte, *http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/apps/%s/sso/saml/metadata?kid=%s", m.baseURL, id, keyID)
	return m.GetXml(url)
}

func (m *ApiSupplement) GetSAMLIdpMetdata(id string) ([]byte, *http.Response, error) {
	url := fmt.Sprintf("%s/api/v1/idps/%s/metadata.xml", m.baseURL, id)
	return m.GetXml(url)
}

func (m *ApiSupplement) GetXml(url string) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", m.token))
	req.Header.Add("User-Agent", "Terraform Okta Provider")
	req.Header.Add("Accept", "application/xml")
	res, err := m.requestExecutor.DoWithRetries(req, 0)
	if err != nil {
		return nil, res, err
	} else if res.StatusCode == http.StatusNotFound {
		return nil, res, nil
	} else if res.StatusCode != http.StatusOK {
		return nil, res, fmt.Errorf("failed to get metadata for url: %s, status: %s", url, res.Status)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)

	return data, res, err
}
