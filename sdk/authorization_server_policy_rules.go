package sdk

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type (
	AuthorizationServerPolicyRule struct {
		Status     string                                   `json:"status,omitempty"`
		Priority   int                                      `json:"priority,omitempty"`
		Type       string                                   `json:"type,omitempty"`
		Name       string                                   `json:"name,omitempty"`
		Id         string                                   `json:"id,omitempty"`
		Conditions *AuthorizationServerPolicyRuleConditions `json:"conditions,omitempty"`
		Actions    *AuthorizationServerPolicyRuleActions    `json:"actions,omitempty"`
	}

	AuthServerInlineHook struct {
		Id string `json:"id,omitempty"`
	}

	AuthorizationServerPolicyRuleActions struct {
		Token *TokenActions `json:"token,omitempty"`
	}

	AuthorizationServerPolicyRuleConditions struct {
		GrantTypes *Whitelist                     `json:"grantTypes,omitempty"`
		People     *okta.GroupRulePeopleCondition `json:"people,omitempty"`
		Scopes     *Whitelist                     `json:"scopes,omitempty"`
	}

	TokenActions struct {
		AccessTokenLifetimeMinutes  int                   `json:"accessTokenLifetimeMinutes,omitempty"`
		RefreshTokenLifetimeMinutes int                   `json:"refreshTokenLifetimeMinutes,omitempty"`
		RefreshTokenWindowMinutes   int                   `json:"refreshTokenWindowMinutes,omitempty"`
		InlineHook                  *AuthServerInlineHook `json:"inlineHook,omitempty"`
	}
)

func (m *ApiSupplement) DeleteAuthorizationServerPolicyRule(ctx context.Context, authServerID, policyID, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/policies/%s/rules/%s", authServerID, policyID, id)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *ApiSupplement) CreateAuthorizationServerPolicyRule(ctx context.Context, authServerID, policyID string, body AuthorizationServerPolicyRule, qp *query.Params) (*AuthorizationServerPolicyRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/policies/%s/rules", authServerID, policyID)
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

func (m *ApiSupplement) UpdateAuthorizationServerPolicyRule(ctx context.Context, authServerID, policyID, id string, body AuthorizationServerPolicyRule, qp *query.Params) (*AuthorizationServerPolicyRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/policies/%s/rules/%s", authServerID, policyID, id)
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

func (m *ApiSupplement) GetAuthorizationServerPolicyRule(ctx context.Context, authServerID, policyID, id string, authorizationServerInstance AuthorizationServerPolicyRule) (*AuthorizationServerPolicyRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/policies/%s/rules/%s", authServerID, policyID, id)
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
