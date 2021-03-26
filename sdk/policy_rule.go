package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

func PasswordPolicyRule() PolicyRule {
	return PolicyRule{Type: PasswordPolicyType}
}

func SignOnPolicyRule() PolicyRule {
	return PolicyRule{Type: SignOnPolicyRuleType}
}

func MfaPolicyRule() PolicyRule {
	return PolicyRule{Type: MfaPolicyType}
}

type PolicyRule struct {
	Id          string                     `json:"id,omitempty"`
	Type        string                     `json:"type,omitempty"`
	Name        string                     `json:"name,omitempty"`
	Status      string                     `json:"status,omitempty"`
	Priority    int64                      `json:"priority,omitempty"`
	System      *bool                      `json:"system,omitempty"`
	Created     *time.Time                 `json:"created,omitempty"`
	LastUpdated *time.Time                 `json:"lastUpdated,omitempty"`
	Conditions  *okta.PolicyRuleConditions `json:"conditions,omitempty"`
	Actions     PolicyRuleActions          `json:"actions,omitempty"`
}

type PolicyRuleActions struct {
	Enroll *Enroll `json:"enroll,omitempty"`
	*okta.OktaSignOnPolicyRuleActions
	*okta.PasswordPolicyRuleActions
}

// Enumerates all policy rules.
func (m *ApiSupplement) ListPolicyRules(ctx context.Context, policyID string) ([]PolicyRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules", policyID)

	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var policyRule []PolicyRule

	resp, err := m.RequestExecutor.Do(ctx, req, &policyRule)
	if err != nil {
		return nil, resp, err
	}

	return policyRule, resp, nil
}

// Creates a policy rule.
func (m *ApiSupplement) CreatePolicyRule(ctx context.Context, policyID string, body PolicyRule) (*PolicyRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules", policyID)

	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var policyRule PolicyRule

	resp, err := m.RequestExecutor.Do(ctx, req, &policyRule)
	if err != nil {
		return nil, resp, err
	}

	return &policyRule, resp, nil
}

// Gets a policy rule.
func (m *ApiSupplement) GetPolicyRule(ctx context.Context, policyID, ruleId string) (*PolicyRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules/%v", policyID, ruleId)

	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var policyRule PolicyRule

	resp, err := m.RequestExecutor.Do(ctx, req, &policyRule)
	if err != nil {
		return nil, resp, err
	}

	return &policyRule, resp, nil
}

// Updates a policy rule.
func (m *ApiSupplement) UpdatePolicyRule(ctx context.Context, policyID, ruleId string, body PolicyRule) (*PolicyRule, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules/%v", policyID, ruleId)

	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var policyRule PolicyRule

	resp, err := m.RequestExecutor.Do(ctx, req, &policyRule)
	if err != nil {
		return nil, resp, err
	}

	return &policyRule, resp, nil
}
