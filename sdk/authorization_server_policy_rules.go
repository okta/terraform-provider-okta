package sdk

import (
	"fmt"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

type (
	AuthorizationServerPolicyRule struct {
		Status     string                `json:"status,omitempty"`
		Priority   int                   `json:"priority,omitempty"`
		Type       string                `json:"type,omitempty"`
		Name       string                `json:"name,omitempty"`
		Id         string                `json:"id,omitempty"`
		Conditions *PolicyRuleConditions `json:"conditions,omitempty"`
		Actions    *PolicyRuleActions    `json:"actions,omitempty"`
	}

	AuthServerInlineHook struct {
		Id string `json:"id,omitempty"`
	}

	PolicyRuleActions struct {
		Token *TokenActions `json:"token,omitempty"`
	}

	PolicyRuleConditions struct {
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

func (m *ApiSupplement) DeleteAuthorizationServerPolicyRule(authServerId, policyId, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/policies/%s/rules/%s", authServerId, policyId, id)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(req, nil)
}

func (m *ApiSupplement) ListAuthorizationServerPolicyRules(authServerId, policyId string) ([]*AuthorizationServerPolicyRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/policies/%s/rules", authServerId, policyId)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var auth []*AuthorizationServerPolicyRule
	resp, err := m.RequestExecutor.Do(req, &auth)
	return auth, resp, err
}

func (m *ApiSupplement) CreateAuthorizationServerPolicyRule(authServerId, policyId string, body AuthorizationServerPolicyRule, qp *query.Params) (*AuthorizationServerPolicyRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/policies/%s/rules", authServerId, policyId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	authorizationServer := body
	resp, err := m.RequestExecutor.Do(req, &authorizationServer)
	return &authorizationServer, resp, err
}

func (m *ApiSupplement) UpdateAuthorizationServerPolicyRule(authServerId, policyId, id string, body AuthorizationServerPolicyRule, qp *query.Params) (*AuthorizationServerPolicyRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/policies/%s/rules/%s", authServerId, policyId, id)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	authorizationServer := body
	resp, err := m.RequestExecutor.Do(req, &authorizationServer)
	if err != nil {
		return nil, resp, err
	}
	return &authorizationServer, resp, nil
}

func (m *ApiSupplement) GetAuthorizationServerPolicyRule(authServerId, policyId, id string, authorizationServerInstance AuthorizationServerPolicyRule) (*AuthorizationServerPolicyRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%s/policies/%s/rules/%s", authServerId, policyId, id)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	authorizationServer := authorizationServerInstance
	resp, err := m.RequestExecutor.Do(req, &authorizationServer)
	if err != nil {
		return nil, resp, err
	}
	return &authorizationServer, resp, nil
}
