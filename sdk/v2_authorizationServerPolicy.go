// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type AuthorizationServerPolicyResource resource

type AuthorizationServerPolicy struct {
	Embedded    interface{}           `json:"_embedded,omitempty"`
	Links       interface{}           `json:"_links,omitempty"`
	Conditions  *PolicyRuleConditions `json:"conditions,omitempty"`
	Created     *time.Time            `json:"created,omitempty"`
	Description string                `json:"description,omitempty"`
	Id          string                `json:"id,omitempty"`
	LastUpdated *time.Time            `json:"lastUpdated,omitempty"`
	Name        string                `json:"name,omitempty"`
	Priority    int64                 `json:"-"`
	PriorityPtr *int64                `json:"priority,omitempty"`
	Status      string                `json:"status,omitempty"`
	System      *bool                 `json:"system,omitempty"`
	Type        string                `json:"type,omitempty"`
}

func (m *AuthorizationServerPolicyResource) GetAuthorizationServerPolicy(ctx context.Context, authServerId, policyId string) (*AuthorizationServerPolicy, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/policies/%v", authServerId, policyId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var authorizationServerPolicy *AuthorizationServerPolicy

	resp, err := rq.Do(ctx, req, &authorizationServerPolicy)
	if err != nil {
		return nil, resp, err
	}

	return authorizationServerPolicy, resp, nil
}

func (m *AuthorizationServerPolicyResource) UpdateAuthorizationServerPolicy(ctx context.Context, authServerId, policyId string, body AuthorizationServerPolicy) (*AuthorizationServerPolicy, *Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/policies/%v", authServerId, policyId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var authorizationServerPolicy *AuthorizationServerPolicy

	resp, err := rq.Do(ctx, req, &authorizationServerPolicy)
	if err != nil {
		return nil, resp, err
	}

	return authorizationServerPolicy, resp, nil
}

func (m *AuthorizationServerPolicyResource) DeleteAuthorizationServerPolicy(ctx context.Context, authServerId, policyId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/authorizationServers/%v/policies/%v", authServerId, policyId)

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

func (a *AuthorizationServerPolicy) MarshalJSON() ([]byte, error) {
	type Alias AuthorizationServerPolicy
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.Priority != 0 {
		result.PriorityPtr = Int64Ptr(a.Priority)
	}
	return json.Marshal(&result)
}

func (a *AuthorizationServerPolicy) UnmarshalJSON(data []byte) error {
	type Alias AuthorizationServerPolicy

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
