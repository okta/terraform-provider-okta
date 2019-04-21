package okta

import (
	"fmt"
	"github.com/okta/okta-sdk-golang/okta"
	"io/ioutil"
	"net/http"
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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("SSWS %s", m.token))
	req.Header.Add("User-Agent", "Terraform Okta Provider")
	req.Header.Add("Accept", "application/xml")
	res, err := m.client.Do(req)
	if err != nil {
		return nil, res, err
	} else if res.StatusCode == http.StatusNotFound {
		return nil, res, nil
	} else if res.StatusCode != http.StatusOK {
		return nil, res, fmt.Errorf("failed to get metadata for app ID: %s, key ID: %s, status: %s", id, keyID, res.Status)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)

	return data, res, err
}
