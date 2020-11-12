package sdk

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type (
	IdpDiscoveryRuleActions struct {
		IDP *IdpDiscoveryRuleIdp `json:"idp"`
	}

	IdpDiscoveryRuleApp struct {
		Exclude []*IdpDiscoveryRuleAppObj `json:"exclude"`
		Include []*IdpDiscoveryRuleAppObj `json:"include"`
	}

	IdpDiscoveryRuleAppObj struct {
		Type string `json:"type,omitempty"`
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	IdpDiscoveryRuleConditions struct {
		App            *IdpDiscoveryRuleApp            `json:"app"`
		Network        *IdpDiscoveryRuleNetwork        `json:"network"`
		Platform       *IdpDiscoveryRulePlatform       `json:"platform,omitempty"`
		UserIdentifier *IdpDiscoveryRuleUserIdentifier `json:"userIdentifier,omitempty"`
	}

	IdpDiscoveryRuleIdp struct {
		Providers []*IdpDiscoveryRuleProvider `json:"providers"`
	}

	IdpDiscoveryRuleNetwork struct {
		Connection string   `json:"connection,omitempty"`
		Include    []string `json:"include,omitempty"`
		Exclude    []string `json:"exclude,omitempty"`
	}

	IdpDiscoveryRulePattern struct {
		MatchType string `json:"matchType,omitempty"`
		Value     string `json:"value,omitempty"`
	}

	IdpDiscoveryRulePlatformOS struct {
		Type       string `json:"type,omitempty"`
		Expression string `json:"expression,omitempty"`
	}

	IdpDiscoveryRulePlatformInclude struct {
		Os   *IdpDiscoveryRulePlatformOS `json:"os"`
		Type string                      `json:"type,omitempty"`
	}

	IdpDiscoveryRulePlatform struct {
		Exclude []interface{}                      `json:"exclude,omitempty"`
		Include []*IdpDiscoveryRulePlatformInclude `json:"include,omitempty"`
	}

	IdpDiscoveryRuleProvider struct {
		Type string `json:"type,omitempty"`
		ID   string `json:"id,omitempty"`
	}

	IdpDiscoveryRuleUserIdentifier struct {
		Attribute string                     `json:"attribute,omitempty"`
		Patterns  []*IdpDiscoveryRulePattern `json:"patterns,omitempty"`
		Type      string                     `json:"type,omitempty"`
	}

	IdpDiscoveryRulePolicy struct {
		Conditions  interface{} `json:"conditions"`
		Created     string      `json:"created"`
		Description string      `json:"description"`
		ID          string      `json:"id"`
		LastUpdated string      `json:"lastUpdated"`
		Name        string      `json:"name"`
		Priority    int64       `json:"priority"`
		Status      string      `json:"status"`
		System      bool        `json:"system"`
		Type        string      `json:"type"`
	}

	IdpDiscoveryRule struct {
		Actions     *IdpDiscoveryRuleActions    `json:"actions"`
		Conditions  *IdpDiscoveryRuleConditions `json:"conditions"`
		Created     string                      `json:"created"`
		ID          string                      `json:"id"`
		LastUpdated string                      `json:"lastUpdated"`
		Name        string                      `json:"name"`
		Priority    int                         `json:"priority"`
		Status      string                      `json:"status"`
		System      bool                        `json:"system"`
		Type        string                      `json:"type"`
	}
)

func (m *ApiSupplement) DeleteIdpDiscoveryRule(policyId, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%s/rules/%s", policyId, id)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(context.Background(), req, nil)
}

func (m *ApiSupplement) ListIdpDiscoveryRules(policyId string) ([]*IdpDiscoveryRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%s/rules", policyId)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var auth []*IdpDiscoveryRule
	resp, err := m.RequestExecutor.Do(context.Background(), req, &auth)
	return auth, resp, err
}

func (m *ApiSupplement) CreateIdpDiscoveryRule(policyId string, body IdpDiscoveryRule, qp *query.Params) (*IdpDiscoveryRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%s/rules", policyId)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	rule := body
	resp, err := m.RequestExecutor.Do(context.Background(), req, &rule)
	return &rule, resp, err
}

func (m *ApiSupplement) UpdateIdpDiscoveryRule(policyId, id string, body IdpDiscoveryRule, qp *query.Params) (*IdpDiscoveryRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%s/rules/%s", policyId, id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	rule := body
	resp, err := m.RequestExecutor.Do(context.Background(), req, &rule)
	return &rule, resp, err
}

func (m *ApiSupplement) GetIdpDiscoveryRule(policyId, id string) (*IdpDiscoveryRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%s/rules/%s", policyId, id)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	rule := &IdpDiscoveryRule{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, rule)
	return rule, resp, err
}

func (m *ApiSupplement) ActivateRule(policyId, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%s/rules/%s/lifecycle/activate", policyId, id)
	req, err := m.RequestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(context.Background(), req, nil)
}

func (m *ApiSupplement) DeactivateRule(policyId, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%s/rules/%s/lifecycle/deactivate", policyId, id)
	req, err := m.RequestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(context.Background(), req, nil)
}
