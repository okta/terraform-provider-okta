// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/okta/terraform-provider-okta/sdk/query"
)

type GroupRuleResource resource

type GroupRule struct {
	Actions     *GroupRuleAction     `json:"actions,omitempty"`
	Conditions  *GroupRuleConditions `json:"conditions,omitempty"`
	Created     *time.Time           `json:"created,omitempty"`
	Id          string               `json:"id,omitempty"`
	LastUpdated *time.Time           `json:"lastUpdated,omitempty"`
	Name        string               `json:"name,omitempty"`
	Status      string               `json:"status,omitempty"`
	Type        string               `json:"type,omitempty"`
}

// Updates a group rule. Only &#x60;INACTIVE&#x60; rules can be updated.
func (m *GroupRuleResource) UpdateGroupRule(ctx context.Context, ruleId string, body GroupRule) (*GroupRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/groups/rules/%v", ruleId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var groupRule *GroupRule

	resp, err := rq.Do(ctx, req, &groupRule)
	if err != nil {
		return nil, resp, err
	}

	return groupRule, resp, nil
}

// Removes a specific group rule by id from your organization
func (m *GroupRuleResource) DeleteGroupRule(ctx context.Context, ruleId string, qp *query.Params) (*Response, error) {
	url := fmt.Sprintf("/api/v1/groups/rules/%v", ruleId)
	if qp != nil {
		url = url + qp.String()
	}

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
