package sdk

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func PasswordPolicyRule() SdkPolicyRule {
	return SdkPolicyRule{Type: PasswordPolicyType}
}

func SignOnPolicyRule() SdkPolicyRule {
	return SdkPolicyRule{Type: SignOnPolicyRuleType}
}

func MfaPolicyRule() SdkPolicyRule {
	return SdkPolicyRule{Type: MfaPolicyType}
}

func ProfileEnrollmentPolicyRule() SdkPolicyRule {
	return SdkPolicyRule{Type: ProfileEnrollmentPolicyType}
}

type SdkPolicyRule struct {
	Id          string                `json:"id,omitempty"`
	Type        string                `json:"type,omitempty"`
	Name        string                `json:"name,omitempty"`
	Status      string                `json:"status,omitempty"`
	Priority    int64                 `json:"priority,omitempty"`
	System      *bool                 `json:"system,omitempty"`
	Created     *time.Time            `json:"created,omitempty"`
	LastUpdated *time.Time            `json:"lastUpdated,omitempty"`
	Conditions  *PolicyRuleConditions `json:"conditions"`
	Actions     SdkPolicyRuleActions  `json:"actions,omitempty"`
}

type SdkPolicyRuleActions struct {
	SignOn            *SdkSignOnPolicyRuleSignOnActions  `json:"signon,omitempty"`
	ProfileEnrollment *ProfileEnrollmentPolicyRuleAction `json:"profileEnrollment,omitempty"`
	*PasswordPolicyRuleActions
}

type SdkSignOnPolicyRuleSignOnActions struct {
	Access                  string                                     `json:"access,omitempty"`
	FactorLifetime          int64                                      `json:"factorLifetime,omitempty"`
	FactorPromptMode        string                                     `json:"factorPromptMode,omitempty"`
	PrimaryFactor           string                                     `json:"primaryFactor,omitempty"`
	RememberDeviceByDefault *bool                                      `json:"rememberDeviceByDefault,omitempty"`
	RequireFactor           *bool                                      `json:"requireFactor,omitempty"`
	Session                 *OktaSignOnPolicyRuleSignonSessionActions  `json:"session,omitempty"`
	Challenge               *SdkSignOnPolicyRuleSignOnActionsChallenge `json:"challenge,omitempty"`
}

type SdkSignOnPolicyRuleSignOnActionsChallenge struct {
	Chain []SdkSignOnPolicyRuleSignOnActionsChallengeChain `json:"chain,omitempty"`
}

type SdkSignOnPolicyRuleSignOnActionsChallengeChain struct {
	Criteria []SdkSignOnPolicyRuleSignOnActionsChallengeChainCriteria `json:"criteria,omitempty"`
	Next     []SdkSignOnPolicyRuleSignOnActionsChallengeChainNext     `json:"next,omitempty"`
}

type SdkSignOnPolicyRuleSignOnActionsChallengeChainCriteria struct {
	Provider   string `json:"provider,omitempty"`
	FactorType string `json:"factorType,omitempty"`
}

type SdkSignOnPolicyRuleSignOnActionsChallengeChainNext struct {
	Criteria []SdkSignOnPolicyRuleSignOnActionsChallengeChainCriteria `json:"criteria,omitempty"`
}

// ListPolicyRules enumerates all policy rules.
func (m *APISupplement) ListPolicyRules(ctx context.Context, policyID string) ([]SdkPolicyRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules", policyID)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var policyRule []SdkPolicyRule
	resp, err := m.RequestExecutor.Do(ctx, req, &policyRule)
	if err != nil {
		return nil, resp, err
	}
	return policyRule, resp, nil
}

// CreatePolicyRule creates a policy rule.
func (m *APISupplement) CreatePolicyRule(ctx context.Context, policyID string, body SdkPolicyRule) (*SdkPolicyRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules", policyID)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var policyRule *SdkPolicyRule
	resp, err := m.RequestExecutor.Do(ctx, req, &policyRule)
	if err != nil {
		return nil, resp, err
	}
	return policyRule, resp, nil
}

// GetPolicyRule gets a policy rule.
func (m *APISupplement) GetPolicyRule(ctx context.Context, policyID, ruleId string) (*SdkPolicyRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules/%v", policyID, ruleId)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var policyRule *SdkPolicyRule
	resp, err := m.RequestExecutor.Do(ctx, req, &policyRule)
	if err != nil {
		return nil, resp, err
	}
	return policyRule, resp, nil
}

// UpdatePolicyRule updates a policy rule.
func (m *APISupplement) UpdatePolicyRule(ctx context.Context, policyID, ruleId string, body SdkPolicyRule) (*SdkPolicyRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v/rules/%v", policyID, ruleId)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}
	var policyRule *SdkPolicyRule
	resp, err := m.RequestExecutor.Do(ctx, req, &policyRule)
	if err != nil {
		return nil, resp, err
	}
	return policyRule, resp, nil
}
