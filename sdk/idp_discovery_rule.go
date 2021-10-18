package sdk

import (
	"context"
	"fmt"
	"net/http"

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

func (m *APISupplement) CreateIdpDiscoveryRule(ctx context.Context, policyID string, body IdpDiscoveryRule, qp *query.Params) (*IdpDiscoveryRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%s/rules", policyID)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	rule := body
	resp, err := m.RequestExecutor.Do(ctx, req, &rule)
	if err != nil {
		return nil, resp, err
	}
	return &rule, resp, err
}

func (m *APISupplement) UpdateIdpDiscoveryRule(ctx context.Context, policyID, id string, body IdpDiscoveryRule, qp *query.Params) (*IdpDiscoveryRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%s/rules/%s", policyID, id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}
	rule := body
	resp, err := m.RequestExecutor.Do(ctx, req, &rule)
	if err != nil {
		return nil, resp, err
	}
	return &rule, resp, err
}

func (m *APISupplement) GetIdpDiscoveryRule(ctx context.Context, policyID, id string) (*IdpDiscoveryRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%s/rules/%s", policyID, id)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var rule *IdpDiscoveryRule
	resp, err := m.RequestExecutor.Do(ctx, req, &rule)
	if err != nil {
		return nil, resp, err
	}
	return rule, resp, nil
}
