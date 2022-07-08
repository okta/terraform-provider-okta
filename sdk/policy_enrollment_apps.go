package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta/query"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

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
