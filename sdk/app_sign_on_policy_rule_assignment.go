package sdk

import (
	"context"
	"fmt"
	"path"
)

type AppSignOnAssignment struct {
	AppId    string
	PolicyId string
}

type appWithAccessPolicy struct {
	Id    string `json:"id,omitempty"`
	Links struct {
		AccessPolicy struct {
			Href string `json:"href,omitempty"`
		} `json:"accessPolicy,omitempty"`
	} `json:"_links,omitempty"`
}

func (m *APISupplement) GetAppSignOnPolicyRuleAssigment(ctx context.Context, appId string) (*AppSignOnAssignment, error) {
	url := fmt.Sprintf("/api/v1/apps/%v", appId)

	req, err := m.RequestExecutor.WithAccept("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var app appWithAccessPolicy
	_, err = m.RequestExecutor.Do(ctx, req, &app)
	if err != nil {
		return nil, err
	}

	policyId := path.Base(app.Links.AccessPolicy.Href)

	return &AppSignOnAssignment{AppId: app.Id, PolicyId: policyId}, nil
}

func (m *APISupplement) SetAppSignOnPolicyRuleAssigment(ctx context.Context, appId string, policyId string) (*AppSignOnAssignment, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/policies/%v", appId, policyId)
	req, err := m.RequestExecutor.WithAccept("application/json").NewRequest("PUT", url, nil)
	if err != nil {
		return nil, err
	}

	_, err = m.RequestExecutor.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	return &AppSignOnAssignment{AppId: appId, PolicyId: policyId}, nil
}
