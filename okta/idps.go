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
func (m *ApiSupplement) ListIdentityProviders(idps interface{}, qp *query.Params) (interface{}, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps")
	if qp != nil {
		url += qp.String()
	}
	req, err := m.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := m.requestExecutor.Do(req, idps)
	return idps, resp, err
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

	resp, err := m.requestExecutor.Do(req, body)
	return body, resp, err
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
	resp, err := m.requestExecutor.Do(req, idp)
	return idp, resp, err
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

func (m *ApiSupplement) GenerateIdentityProviderSigningKey(idpId string, yearsValid int) (*SigningKey, *okta.Response, error) {
	key := &SigningKey{}
	url := fmt.Sprintf("/api/v1/idps/%s/credentials/keys/generate?validityYears=%d", idpId, yearsValid)
	req, err := m.requestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return key, nil, err
	}

	resp, err := m.requestExecutor.Do(req, key)
	return key, resp, err
}

func (m *ApiSupplement) GetIdentityProviderSigningKey(idpId, kid string) (*SigningKey, *okta.Response, error) {
	key := &SigningKey{}
	url := fmt.Sprintf("/api/v1/idps/%s/credentials/keys/%s", idpId, kid)
	req, err := m.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return key, nil, err
	}

	resp, err := m.requestExecutor.Do(req, key)
	return key, resp, err
}

func (m *ApiSupplement) DeleteIdentityProviderSigningKey(idpId, kid string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps/%s/credentials/keys/%s", idpId, kid)
	req, err := m.requestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.requestExecutor.Do(req, nil)
}

func (m *ApiSupplement) AddIdentityProviderCertificate(cert *Certificate) (*SigningKey, *okta.Response, error) {
	key := &SigningKey{}
	req, err := m.requestExecutor.NewRequest("POST", "/api/v1/idps/credentials/keys", cert)
	if err != nil {
		return key, nil, err
	}
	resp, err := m.requestExecutor.Do(req, key)
	return key, resp, err
}

func (m *ApiSupplement) GetIdentityProviderCertificate(kid string) (*SigningKey, *okta.Response, error) {
	key := &SigningKey{}
	url := fmt.Sprintf("/api/v1/idps/credentials/keys/%s", kid)
	req, err := m.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return key, nil, err
	}

	resp, err := m.requestExecutor.Do(req, key)
	return key, resp, err
}

func (m *ApiSupplement) DeleteIdentityProviderCertificate(kid string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps/credentials/keys/%s", kid)
	req, err := m.requestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.requestExecutor.Do(req, nil)
	return resp, err
}
