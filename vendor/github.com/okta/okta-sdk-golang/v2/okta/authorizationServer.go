/*
* Copyright 2018 - Present Okta, Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

// AUTO-GENERATED!  DO NOT EDIT FILE DIRECTLY

package okta

import (
	"context"
	"fmt"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"time"
)

type AuthorizationServerResource resource

type AuthorizationServer struct {
	Links       interface{}                     `json:"_links,omitempty"`
	Audiences   []string                        `json:"audiences,omitempty"`
	Created     *time.Time                      `json:"created,omitempty"`
	Credentials *AuthorizationServerCredentials `json:"credentials,omitempty"`
	Description string                          `json:"description,omitempty"`
	Id          string                          `json:"id,omitempty"`
	Issuer      string                          `json:"issuer,omitempty"`
	IssuerMode  string                          `json:"issuerMode,omitempty"`
	LastUpdated *time.Time                      `json:"lastUpdated,omitempty"`
	Name        string                          `json:"name,omitempty"`
	Status      string                          `json:"status,omitempty"`
}

func (m *AuthorizationServerResource) CreateAuthorizationServer(ctx context.Context, body AuthorizationServer) (*AuthorizationServer, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers")

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var authorizationServer *AuthorizationServer

	resp, err := m.client.requestExecutor.Do(ctx, req, &authorizationServer)
	if err != nil {
		return nil, resp, err
	}

	return authorizationServer, resp, nil
}

func (m *AuthorizationServerResource) GetAuthorizationServer(ctx context.Context, authServerId string) (*AuthorizationServer, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v", authServerId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var authorizationServer *AuthorizationServer

	resp, err := m.client.requestExecutor.Do(ctx, req, &authorizationServer)
	if err != nil {
		return nil, resp, err
	}

	return authorizationServer, resp, nil
}

func (m *AuthorizationServerResource) UpdateAuthorizationServer(ctx context.Context, authServerId string, body AuthorizationServer) (*AuthorizationServer, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v", authServerId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var authorizationServer *AuthorizationServer

	resp, err := m.client.requestExecutor.Do(ctx, req, &authorizationServer)
	if err != nil {
		return nil, resp, err
	}

	return authorizationServer, resp, nil
}

func (m *AuthorizationServerResource) DeleteAuthorizationServer(ctx context.Context, authServerId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v", authServerId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (m *AuthorizationServerResource) ListAuthorizationServers(ctx context.Context, qp *query.Params) ([]*AuthorizationServer, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers")
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var authorizationServer []*AuthorizationServer

	resp, err := m.client.requestExecutor.Do(ctx, req, &authorizationServer)
	if err != nil {
		return nil, resp, err
	}

	return authorizationServer, resp, nil
}

func (m *AuthorizationServerResource) ListOAuth2Claims(ctx context.Context, authServerId string) ([]*OAuth2Claim, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/claims", authServerId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2Claim []*OAuth2Claim

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2Claim)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2Claim, resp, nil
}

func (m *AuthorizationServerResource) CreateOAuth2Claim(ctx context.Context, authServerId string, body OAuth2Claim) (*OAuth2Claim, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/claims", authServerId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2Claim *OAuth2Claim

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2Claim)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2Claim, resp, nil
}

func (m *AuthorizationServerResource) DeleteOAuth2Claim(ctx context.Context, authServerId string, claimId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/claims/%v", authServerId, claimId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (m *AuthorizationServerResource) GetOAuth2Claim(ctx context.Context, authServerId string, claimId string) (*OAuth2Claim, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/claims/%v", authServerId, claimId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2Claim *OAuth2Claim

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2Claim)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2Claim, resp, nil
}

func (m *AuthorizationServerResource) UpdateOAuth2Claim(ctx context.Context, authServerId string, claimId string, body OAuth2Claim) (*OAuth2Claim, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/claims/%v", authServerId, claimId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2Claim *OAuth2Claim

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2Claim)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2Claim, resp, nil
}

func (m *AuthorizationServerResource) ListOAuth2ClientsForAuthorizationServer(ctx context.Context, authServerId string) ([]*OAuth2Client, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/clients", authServerId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2Client []*OAuth2Client

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2Client)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2Client, resp, nil
}

func (m *AuthorizationServerResource) RevokeRefreshTokensForAuthorizationServerAndClient(ctx context.Context, authServerId string, clientId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/clients/%v/tokens", authServerId, clientId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (m *AuthorizationServerResource) ListRefreshTokensForAuthorizationServerAndClient(ctx context.Context, authServerId string, clientId string, qp *query.Params) ([]*OAuth2RefreshToken, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/clients/%v/tokens", authServerId, clientId)
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2RefreshToken []*OAuth2RefreshToken

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2RefreshToken)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2RefreshToken, resp, nil
}

func (m *AuthorizationServerResource) RevokeRefreshTokenForAuthorizationServerAndClient(ctx context.Context, authServerId string, clientId string, tokenId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/clients/%v/tokens/%v", authServerId, clientId, tokenId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (m *AuthorizationServerResource) GetRefreshTokenForAuthorizationServerAndClient(ctx context.Context, authServerId string, clientId string, tokenId string, qp *query.Params) (*OAuth2RefreshToken, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/clients/%v/tokens/%v", authServerId, clientId, tokenId)
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2RefreshToken *OAuth2RefreshToken

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2RefreshToken)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2RefreshToken, resp, nil
}

func (m *AuthorizationServerResource) ListAuthorizationServerKeys(ctx context.Context, authServerId string) ([]*JsonWebKey, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/credentials/keys", authServerId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var jsonWebKey []*JsonWebKey

	resp, err := m.client.requestExecutor.Do(ctx, req, &jsonWebKey)
	if err != nil {
		return nil, resp, err
	}

	return jsonWebKey, resp, nil
}

func (m *AuthorizationServerResource) RotateAuthorizationServerKeys(ctx context.Context, authServerId string, body JwkUse) ([]*JsonWebKey, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/credentials/lifecycle/keyRotate", authServerId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var jsonWebKey []*JsonWebKey

	resp, err := m.client.requestExecutor.Do(ctx, req, &jsonWebKey)
	if err != nil {
		return nil, resp, err
	}

	return jsonWebKey, resp, nil
}

func (m *AuthorizationServerResource) ActivateAuthorizationServer(ctx context.Context, authServerId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/lifecycle/activate", authServerId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (m *AuthorizationServerResource) DeactivateAuthorizationServer(ctx context.Context, authServerId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/lifecycle/deactivate", authServerId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (m *AuthorizationServerResource) ListAuthorizationServerPolicies(ctx context.Context, authServerId string) ([]*Policy, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/policies", authServerId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var policy []*Policy

	resp, err := m.client.requestExecutor.Do(ctx, req, &policy)
	if err != nil {
		return nil, resp, err
	}

	return policy, resp, nil
}

func (m *AuthorizationServerResource) CreateAuthorizationServerPolicy(ctx context.Context, authServerId string, body Policy) (*Policy, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/policies", authServerId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var policy *Policy

	resp, err := m.client.requestExecutor.Do(ctx, req, &policy)
	if err != nil {
		return nil, resp, err
	}

	return policy, resp, nil
}

func (m *AuthorizationServerResource) DeleteAuthorizationServerPolicy(ctx context.Context, authServerId string, policyId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/policies/%v", authServerId, policyId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (m *AuthorizationServerResource) GetAuthorizationServerPolicy(ctx context.Context, authServerId string, policyId string) (*Policy, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/policies/%v", authServerId, policyId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var policy *Policy

	resp, err := m.client.requestExecutor.Do(ctx, req, &policy)
	if err != nil {
		return nil, resp, err
	}

	return policy, resp, nil
}

func (m *AuthorizationServerResource) UpdateAuthorizationServerPolicy(ctx context.Context, authServerId string, policyId string, body Policy) (*Policy, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/policies/%v", authServerId, policyId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var policy *Policy

	resp, err := m.client.requestExecutor.Do(ctx, req, &policy)
	if err != nil {
		return nil, resp, err
	}

	return policy, resp, nil
}

func (m *AuthorizationServerResource) ListOAuth2Scopes(ctx context.Context, authServerId string, qp *query.Params) ([]*OAuth2Scope, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/scopes", authServerId)
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2Scope []*OAuth2Scope

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2Scope)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2Scope, resp, nil
}

func (m *AuthorizationServerResource) CreateOAuth2Scope(ctx context.Context, authServerId string, body OAuth2Scope) (*OAuth2Scope, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/scopes", authServerId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2Scope *OAuth2Scope

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2Scope)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2Scope, resp, nil
}

func (m *AuthorizationServerResource) DeleteOAuth2Scope(ctx context.Context, authServerId string, scopeId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/scopes/%v", authServerId, scopeId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (m *AuthorizationServerResource) GetOAuth2Scope(ctx context.Context, authServerId string, scopeId string) (*OAuth2Scope, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/scopes/%v", authServerId, scopeId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2Scope *OAuth2Scope

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2Scope)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2Scope, resp, nil
}

func (m *AuthorizationServerResource) UpdateOAuth2Scope(ctx context.Context, authServerId string, scopeId string, body OAuth2Scope) (*OAuth2Scope, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/scopes/%v", authServerId, scopeId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2Scope *OAuth2Scope

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2Scope)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2Scope, resp, nil
}
