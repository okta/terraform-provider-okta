package okta

// Not all APIs are supported by okta-sdk-golang, this is one

import (
	"fmt"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func (m *ApiSupplement) DeleteIdentityProvider(id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps/%s", id)
	req, err := m.requestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.requestExecutor.Do(req, nil)
}
func (m *ApiSupplement) ListIdentityProviders(idps []IdentityProvider) ([]IdentityProvider, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps")
	req, err := m.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var auth []IdentityProvider
	resp, err := m.requestExecutor.Do(req, &auth)
	return auth, resp, err
}
func (m *ApiSupplement) CreateIdentityProvider(body IdentityProvider, qp *query.Params) (IdentityProvider, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps")
	if qp != nil {
		url += qp.String()
	}
	req, err := m.requestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	identityProvider := body
	resp, err := m.requestExecutor.Do(req, &identityProvider)
	return identityProvider, resp, err
}

func (m *ApiSupplement) UpdateIdentityProvider(id string, body IdentityProvider, qp *query.Params) (IdentityProvider, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps/%s", id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.requestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	identityProvider := body
	resp, err := m.requestExecutor.Do(req, &identityProvider)
	return identityProvider, resp, err
}

func (m *ApiSupplement) GetIdentityProvider(id string, idp IdentityProvider) (IdentityProvider, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps/%s", id)
	req, err := m.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	identityProvider := idp
	resp, err := m.requestExecutor.Do(req, identityProvider)
	return identityProvider, resp, err
}
func (m *ApiSupplement) ActivateIdentityProvider(id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps/%s/lifecycle/activate", id)
	req, err := m.requestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	return m.requestExecutor.Do(req, nil)
}
func (m *ApiSupplement) DeactivateIdentityProvider(id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps/%s/lifecycle/deactivate", id)
	req, err := m.requestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	return m.requestExecutor.Do(req, nil)
}
