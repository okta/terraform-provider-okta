package sdk

// Not all APIs are supported by okta-sdk-golang, this is one

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func (m *ApiSupplement) DeleteIdentityProvider(ctx context.Context, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps/%s", id)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *ApiSupplement) ListIdentityProviders(ctx context.Context, idps interface{}, qp *query.Params) (interface{}, *okta.Response, error) {
	url := "/api/v1/idps"
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := m.RequestExecutor.Do(ctx, req, idps)
	return idps, resp, err
}

func (m *ApiSupplement) CreateIdentityProvider(ctx context.Context, body IdentityProvider, qp *query.Params) (IdentityProvider, *okta.Response, error) {
	url := "/api/v1/idps"
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	resp, err := m.RequestExecutor.Do(ctx, req, body)
	return body, resp, err
}

func (m *ApiSupplement) UpdateIdentityProvider(ctx context.Context, id string, body IdentityProvider, qp *query.Params) (IdentityProvider, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps/%s", id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	identityProvider := body
	resp, err := m.RequestExecutor.Do(ctx, req, &identityProvider)
	return identityProvider, resp, err
}

func (m *ApiSupplement) GetIdentityProvider(ctx context.Context, id string, idp IdentityProvider) (IdentityProvider, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps/%s", id)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := m.RequestExecutor.Do(ctx, req, idp)
	return idp, resp, err
}

func (m *ApiSupplement) ActivateIdentityProvider(ctx context.Context, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps/%s/lifecycle/activate", id)
	req, err := m.RequestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *ApiSupplement) DeactivateIdentityProvider(ctx context.Context, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps/%s/lifecycle/deactivate", id)
	req, err := m.RequestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *ApiSupplement) GenerateIdentityProviderSigningKey(ctx context.Context, idpId string, yearsValid int) (*SigningKey, *okta.Response, error) {
	key := &SigningKey{}
	url := fmt.Sprintf("/api/v1/idps/%s/credentials/keys/generate?validityYears=%d", idpId, yearsValid)
	req, err := m.RequestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return key, nil, err
	}

	resp, err := m.RequestExecutor.Do(ctx, req, key)
	return key, resp, err
}

func (m *ApiSupplement) GetIdentityProviderSigningKey(ctx context.Context, idpId, kid string) (*SigningKey, *okta.Response, error) {
	key := &SigningKey{}
	url := fmt.Sprintf("/api/v1/idps/%s/credentials/keys/%s", idpId, kid)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return key, nil, err
	}

	resp, err := m.RequestExecutor.Do(ctx, req, key)
	return key, resp, err
}

func (m *ApiSupplement) DeleteIdentityProviderSigningKey(ctx context.Context, kid string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps/credentials/keys/%s", kid)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *ApiSupplement) AddIdentityProviderCertificate(ctx context.Context, cert *Certificate) (*SigningKey, *okta.Response, error) {
	key := &SigningKey{}
	req, err := m.RequestExecutor.NewRequest("POST", "/api/v1/idps/credentials/keys", cert)
	if err != nil {
		return key, nil, err
	}
	resp, err := m.RequestExecutor.Do(ctx, req, key)
	return key, resp, err
}

func (m *ApiSupplement) GetIdentityProviderCertificate(ctx context.Context, kid string) (*SigningKey, *okta.Response, error) {
	key := &SigningKey{}
	url := fmt.Sprintf("/api/v1/idps/credentials/keys/%s", kid)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return key, nil, err
	}

	resp, err := m.RequestExecutor.Do(ctx, req, key)
	return key, resp, err
}

func (m *ApiSupplement) DeleteIdentityProviderCertificate(ctx context.Context, kid string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/idps/credentials/keys/%s", kid)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.RequestExecutor.Do(ctx, req, nil)
	return resp, err
}
