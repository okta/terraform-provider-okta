package sdk

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type AuthorizationServerScope struct {
	Name            string `json:"name,omitempty"`
	Id              string `json:"id,omitempty"`
	Description     string `json:"description,omitempty"`
	Consent         string `json:"consent,omitempty"`
	MetadataPublish string `json:"metadataPublish,omitempty"`
	Default         bool   `json:"default"`
}

func (m *ApiSupplement) DeleteAuthorizationServerScope(ctx context.Context, authServerID, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/scopes/%s", authServerID, id)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *ApiSupplement) CreateAuthorizationServerScope(ctx context.Context, authServerID string, body AuthorizationServerScope, qp *query.Params) (*AuthorizationServerScope, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/scopes", authServerID)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	authorizationServer := body
	resp, err := m.RequestExecutor.Do(ctx, req, &authorizationServer)
	return &authorizationServer, resp, err
}

func (m *ApiSupplement) UpdateAuthorizationServerScope(ctx context.Context, authServerID, id string, body AuthorizationServerScope, qp *query.Params) (*AuthorizationServerScope, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/scopes/%s", authServerID, id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	authorizationServer := body
	resp, err := m.RequestExecutor.Do(ctx, req, &authorizationServer)
	if err != nil {
		return nil, resp, err
	}
	return &authorizationServer, resp, nil
}

func (m *ApiSupplement) GetAuthorizationServerScope(ctx context.Context, authServerID, id string, authorizationServerInstance AuthorizationServerScope) (*AuthorizationServerScope, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/scopes/%s", authServerID, id)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	authorizationServer := authorizationServerInstance
	resp, err := m.RequestExecutor.Do(ctx, req, &authorizationServer)
	if err != nil {
		return nil, resp, err
	}
	return &authorizationServer, resp, nil
}
