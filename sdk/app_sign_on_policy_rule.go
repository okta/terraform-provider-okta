package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

// CreateAppSignOnPolicyRule creates a policy rule.
func (m *APISupplement) CreateAppSignOnPolicyRule(ctx context.Context, policyID string, body okta.AccessPolicyRule) (*okta.AccessPolicyRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules", policyID)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var appSignOnPolicyRule *okta.AccessPolicyRule
	resp, err := m.RequestExecutor.Do(ctx, req, &appSignOnPolicyRule)
	if err != nil {
		return nil, resp, err
	}
	return appSignOnPolicyRule, resp, nil
}

// GetAppSignOnPolicyRule gets a policy rule.
func (m *APISupplement) GetAppSignOnPolicyRule(ctx context.Context, policyID, ruleId string) (*okta.AccessPolicyRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules/%v", policyID, ruleId)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var appSignOnPolicyRule *okta.AccessPolicyRule
	resp, err := m.RequestExecutor.Do(ctx, req, &appSignOnPolicyRule)
	if err != nil {
		return nil, resp, err
	}
	return appSignOnPolicyRule, resp, nil
}

// UpdateAppSignOnPolicyRule updates a policy rule.
func (m *APISupplement) UpdateAppSignOnPolicyRule(ctx context.Context, policyID, ruleId string, body okta.AccessPolicyRule) (*okta.AccessPolicyRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules/%v", policyID, ruleId)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}
	var appSignOnPolicyRule *okta.AccessPolicyRule
	resp, err := m.RequestExecutor.Do(ctx, req, &appSignOnPolicyRule)
	if err != nil {
		return nil, resp, err
	}
	return appSignOnPolicyRule, resp, nil
}

// DeleteAppSignOnPolicyRule deletes app sign on policy rule by ID
func (m *APISupplement) DeleteAppSignOnPolicyRule(ctx context.Context, policyID, ruleId string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules/%v", policyID, ruleId)
	req, err := m.RequestExecutor.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}

// ActivateAppSignOnPolicyRule activates the app sign on policy rule.
func (m *APISupplement) ActivateAppSignOnPolicyRule(ctx context.Context, policyID, ruleID string) (*okta.Response, error) {
	return m.lifecycleChangeAppSignOnPolicyRule(ctx, policyID, ruleID, "activate")
}

// DeactivateAppSignOnPolicyRule deactivates the app sign on policy rule.
func (m *APISupplement) DeactivateAppSignOnPolicyRule(ctx context.Context, policyID, ruleID string) (*okta.Response, error) {
	return m.lifecycleChangeAppSignOnPolicyRule(ctx, policyID, ruleID, "deactivate")
}

func (m *APISupplement) lifecycleChangeAppSignOnPolicyRule(ctx context.Context, policyID, ruleID, action string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%s/rules/%s/lifecycle/%s", policyID, ruleID, action)
	req, err := m.RequestExecutor.
		WithAccept("application/json").
		WithContentType("application/json").
		NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}
