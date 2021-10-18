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
func (m *APISupplement) ListBehaviors(ctx context.Context, qp *query.Params) ([]*Behavior, *okta.Response, error) {
	url := "/api/v1/behaviors"
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
func (m *APISupplement) GetBehavior(ctx context.Context, id string) (*Behavior, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/behaviors/%s", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var behavior *Behavior
	resp, err := m.RequestExecutor.Do(ctx, req, &behavior)
	if err != nil {
		return nil, resp, err
	}
	return behavior, resp, nil
}

// CreateBehavior creates behavior
func (m *APISupplement) CreateBehavior(ctx context.Context, body Behavior) (*Behavior, *okta.Response, error) {
	url := "/api/v1/behaviors"
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var behavior *Behavior
	resp, err := m.RequestExecutor.Do(ctx, req, &behavior)
	if err != nil {
		return nil, resp, err
	}
	return behavior, resp, nil
}

// UpdateBehavior updates behavior
func (m *APISupplement) UpdateBehavior(ctx context.Context, id string, body Behavior) (*Behavior, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/behaviors/%s", id)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}
	var behavior *Behavior
	resp, err := m.RequestExecutor.Do(ctx, req, &behavior)
	if err != nil {
		return nil, resp, err
	}
	return behavior, resp, nil
}

// DeleteBehavior deletes behavior by ID
func (m *APISupplement) DeleteBehavior(ctx context.Context, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/behaviors/%s", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *APISupplement) ActivateBehavior(ctx context.Context, id string) (*okta.Response, error) {
	return m.changeBehaviorLifecycle(ctx, id, "activate")
}

func (m *APISupplement) DeactivateBehavior(ctx context.Context, id string) (*okta.Response, error) {
	return m.changeBehaviorLifecycle(ctx, id, "deactivate")
}

func (m *APISupplement) changeBehaviorLifecycle(ctx context.Context, id, action string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/behaviors/%s/lifecycle/%s", id, action)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := m.RequestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
