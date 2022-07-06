package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta/query"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type AddAppToEnrollmentPolicyRequest struct {
	ResourceType string `json:"resourceType"`
	ResourceId   string `json:"resourceId"`
}

type AddAppToEnrollmentPolicyResponse struct {
	Id    string      `json:"id"`
	Links interface{} `json:"_links,omitempty"`
}

type AddEnrollmentPolicyToAppRequest struct {
	Id string `json:"id"`
}

// AddAppToEnrollmentPolicy adds an app to the policy
func (m *APISupplement) AddAppToEnrollmentPolicy(ctx context.Context, policyID string, body AddAppToEnrollmentPolicyRequest) (*AddAppToEnrollmentPolicyResponse, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%s/policies/%s", body.ResourceId, policyID)
	re := m.cloneRequestExecutor()
	requestBody := &AddEnrollmentPolicyToAppRequest{Id: body.ResourceId}
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, requestBody)
	if err != nil {
		return nil, nil, err
	}
	var response *AddAppToEnrollmentPolicyResponse
	resp, err := re.Do(ctx, req, &response)
	if err != nil {
		return nil, resp, err
	}
	return response, resp, nil
}

func (m *APISupplement) ListEnrollmentPolicyApps(ctx context.Context, policyID string, qp *query.Params) ([]*okta.Application, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%s/app", policyID)
	if qp != nil {
		url += qp.String()
	}
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var applications []*okta.Application
	resp, err := re.Do(ctx, req, &applications)
	if err != nil {
		return nil, resp, err
	}
	return applications, resp, nil
}
