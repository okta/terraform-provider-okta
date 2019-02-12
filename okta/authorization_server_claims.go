package okta

import (
	"fmt"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

type AuthorizationServerClaim struct {
	Status               string      `json:"status,omitempty"`
	ClaimType            string      `json:"claimType,omitempty"`
	ValueType            string      `json:"valueType,omitempty"`
	Value                string      `json:"value,omitempty"`
	AlwaysIncludeInToken bool        `json:"alwaysIncludeInToken,omitempty"`
	Name                 string      `json:"name,omitempty"`
	Id                   string      `json:"id,omitempty"`
	Conditions           *Conditions `json:"conditions"`
}

type Conditions struct {
	Scopes []string `json:"scopes"`
}

func (m *ApiSupplement) DeleteAuthorizationServerClaim(authServerId, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/claims/%s", authServerId, id)
	req, err := m.requestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.requestExecutor.Do(req, nil)
}

func (m *ApiSupplement) ListAuthorizationServerClaims(authServerId string) ([]*AuthorizationServerClaim, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/claims", authServerId)
	req, err := m.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var auth []*AuthorizationServerClaim
	resp, err := m.requestExecutor.Do(req, &auth)
	return auth, resp, err
}

func (m *ApiSupplement) CreateAuthorizationServerClaim(authServerId string, body AuthorizationServerClaim, qp *query.Params) (*AuthorizationServerClaim, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/claims", authServerId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.requestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	authorizationServer := body
	resp, err := m.requestExecutor.Do(req, &authorizationServer)
	return &authorizationServer, resp, err
}

func (m *ApiSupplement) UpdateAuthorizationServerClaim(authServerId, id string, body AuthorizationServerClaim, qp *query.Params) (*AuthorizationServerClaim, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/claims/%s", authServerId, id)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.requestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	authorizationServer := body
	resp, err := m.requestExecutor.Do(req, &authorizationServer)
	if err != nil {
		return nil, resp, err
	}
	return &authorizationServer, resp, nil
}

func (m *ApiSupplement) GetAuthorizationServerClaim(authServerId, id string, authorizationServerInstance AuthorizationServerClaim) (*AuthorizationServerClaim, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/claims/%s", authServerId, id)
	req, err := m.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	authorizationServer := authorizationServerInstance
	resp, err := m.requestExecutor.Do(req, &authorizationServer)
	if err != nil {
		return nil, resp, err
	}
	return &authorizationServer, resp, nil
}
