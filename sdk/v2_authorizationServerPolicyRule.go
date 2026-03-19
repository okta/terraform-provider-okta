// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type AuthorizationServerPolicyRuleResource resource

type AuthorizationServerPolicyRule struct {
	Actions     *AuthorizationServerPolicyRuleActions    `json:"actions,omitempty"`
	Conditions  *AuthorizationServerPolicyRuleConditions `json:"conditions,omitempty"`
	Created     *time.Time                               `json:"created,omitempty"`
	Id          string                                   `json:"id,omitempty"`
	LastUpdated *time.Time                               `json:"lastUpdated,omitempty"`
	Name        string                                   `json:"name,omitempty"`
	Priority    int64                                    `json:"-"`
	PriorityPtr *int64                                   `json:"priority,omitempty"`
	Status      string                                   `json:"status,omitempty"`
	System      *bool                                    `json:"system,omitempty"`
	Type        string                                   `json:"type,omitempty"`
}

// Updates the configuration of the Policy Rule defined in the specified Custom Authorization Server and Policy.
func (m *AuthorizationServerPolicyRuleResource) UpdateAuthorizationServerPolicyRule(ctx context.Context, authServerId, policyId, ruleId string, body AuthorizationServerPolicyRule) (*AuthorizationServerPolicyRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/policies/%v/rules/%v", authServerId, policyId, ruleId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var authorizationServerPolicyRule *AuthorizationServerPolicyRule

	resp, err := rq.Do(ctx, req, &authorizationServerPolicyRule)
	if err != nil {
		return nil, resp, err
	}

	return authorizationServerPolicyRule, resp, nil
}

// Deletes a Policy Rule defined in the specified Custom Authorization Server and Policy.
func (m *AuthorizationServerPolicyRuleResource) DeleteAuthorizationServerPolicyRule(ctx context.Context, authServerId, policyId, ruleId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/policies/%v/rules/%v", authServerId, policyId, ruleId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (a *AuthorizationServerPolicyRule) MarshalJSON() ([]byte, error) {
	type Alias AuthorizationServerPolicyRule
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.Priority != 0 {
		result.PriorityPtr = Int64Ptr(a.Priority)
	}
	return json.Marshal(&result)
}

func (a *AuthorizationServerPolicyRule) UnmarshalJSON(data []byte) error {
	type Alias AuthorizationServerPolicyRule

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.PriorityPtr != nil {
		a.Priority = *result.PriorityPtr
		a.PriorityPtr = result.PriorityPtr
	}
	return nil
}
