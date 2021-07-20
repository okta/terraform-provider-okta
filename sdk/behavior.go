package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type Behavior struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Status   string                 `json:"status"`
	Settings map[string]interface{} `json:"settings"`
	Type     string                 `json:"type"`
}

// ListBehaviors Gets all behaviors based on the query params
func (m *ApiSupplement) ListBehaviors(ctx context.Context, qp *query.Params) ([]*Behavior, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/behaviors")
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var behaviors []*Behavior
	resp, err := m.RequestExecutor.Do(ctx, req, &behaviors)
	if err != nil {
		return nil, resp, err
	}
	return behaviors, resp, nil
}

// GetBehavior gets behavior by ID
func (m *ApiSupplement) GetBehavior(ctx context.Context, id string) (*Behavior, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/behaviors/%s", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var behavior Behavior
	resp, err := m.RequestExecutor.Do(ctx, req, &behavior)
	if err != nil {
		return nil, resp, err
	}
	return &behavior, resp, nil
}
