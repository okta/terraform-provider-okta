package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/terraform-provider-okta/sdk/query"
)

func (m *APISupplement) ListEnrollmentPolicyApps(ctx context.Context, policyID string, qp *query.Params) ([]*Application, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%s/app", policyID)
	if qp != nil {
		url += qp.String()
	}
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var applications []*Application
	resp, err := re.Do(ctx, req, &applications)
	if err != nil {
		return nil, resp, err
	}
	return applications, resp, nil
}
