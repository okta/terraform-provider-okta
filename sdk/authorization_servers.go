package sdk

// Not all APIs are supported by okta-sdk-golang, this is one

import (
	"context"
	"fmt"
	"net/url"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/peterhellberg/link"
)

type AuthorizationServer struct {
	Audiences   []string               `json:"audiences,omitempty"`
	Credentials *AuthServerCredentials `json:"credentials,omitempty"`
	Description string                 `json:"description,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Id          string                 `json:"id,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Issuer      string                 `json:"issuer,omitempty"`
	IssuerMode  string                 `json:"issuerMode,omitempty"`
}

type AuthServerCredentials struct {
	Signing *okta.ApplicationCredentialsSigning `json:"signing,omitempty"`
}

func (m *ApiSupplement) DeleteAuthorizationServer(ctx context.Context, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s", id)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *ApiSupplement) ListAuthorizationServers(ctx context.Context) ([]*AuthorizationServer, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("GET", "/api/v1/authorizationServers", nil)
	if err != nil {
		return nil, nil, err
	}

	var auth []*AuthorizationServer
	resp, err := m.RequestExecutor.Do(ctx, req, &auth)
	return auth, resp, err
}

func (m *ApiSupplement) CreateAuthorizationServer(ctx context.Context, body AuthorizationServer, qp *query.Params) (*AuthorizationServer, *okta.Response, error) {
	url := "/api/v1/authorizationServers"
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

func (m *ApiSupplement) UpdateAuthorizationServer(ctx context.Context, id string, body AuthorizationServer, qp *query.Params) (*AuthorizationServer, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s", id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	authorizationServer := body
	resp, err := m.RequestExecutor.Do(ctx, req, &authorizationServer)
	return &authorizationServer, resp, err
}

func (m *ApiSupplement) GetAuthorizationServer(ctx context.Context, id string) (*AuthorizationServer, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s", id)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	authorizationServer := &AuthorizationServer{}
	resp, err := m.RequestExecutor.Do(ctx, req, authorizationServer)
	return authorizationServer, resp, err
}

func (m *ApiSupplement) ActivateAuthorizationServer(ctx context.Context, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/lifecycle/activate", id)
	req, err := m.RequestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *ApiSupplement) DeactivateAuthorizationServer(ctx context.Context, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/lifecycle/deactivate", id)
	req, err := m.RequestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *ApiSupplement) FindAuthServer(ctx context.Context, name string, qp *query.Params) (*AuthorizationServer, error) {
	authServerList, res, err := m.ListAuthorizationServers(ctx)
	if err != nil {
		return nil, err
	}

	for _, authServer := range authServerList {
		if authServer.Name == name {
			return authServer, nil
		}
	}

	if after := getNextLinkOffset(res); after != "" {
		qp.After = after
		return m.FindAuthServer(ctx, name, qp)
	}
	return nil, nil
}

func (m *ApiSupplement) FilterAuthServers(ctx context.Context, qp *query.Params, arr []*AuthorizationServer, compare func(string) bool) ([]*AuthorizationServer, error) {
	authServerList, res, err := m.ListAuthorizationServers(ctx)
	if err != nil {
		return nil, err
	}

	for _, authServer := range authServerList {
		if compare(authServer.Name) {
			arr = append(arr, authServer)
		}
	}

	if after := getNextLinkOffset(res); after != "" {
		qp.After = after
		return m.FilterAuthServers(ctx, qp, arr, compare)
	}

	return arr, nil
}

func getNextLinkOffset(res *okta.Response) string {
	linkList := link.Parse(res.Header.Get("link"))

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
