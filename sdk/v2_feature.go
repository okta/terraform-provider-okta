package sdk

import (
	"context"
	"fmt"

	"github.com/okta/terraform-provider-okta/sdk/query"
)

type FeatureResource resource

type Feature struct {
	Links       interface{}   `json:"_links,omitempty"`
	Description string        `json:"description,omitempty"`
	Id          string        `json:"id,omitempty"`
	Name        string        `json:"name,omitempty"`
	Stage       *FeatureStage `json:"stage,omitempty"`
	Status      string        `json:"status,omitempty"`
	Type        string        `json:"type,omitempty"`
}

func (m *FeatureResource) GetFeature(ctx context.Context, featureId string) (*Feature, *Response, error) {
	url := fmt.Sprintf("/api/v1/features/%v", featureId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var feature *Feature

	resp, err := rq.Do(ctx, req, &feature)
	if err != nil {
		return nil, resp, err
	}

	return feature, resp, nil
}

func (m *FeatureResource) ListFeatures(ctx context.Context) ([]*Feature, *Response, error) {
	url := "/api/v1/features"

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var feature []*Feature

	resp, err := rq.Do(ctx, req, &feature)
	if err != nil {
		return nil, resp, err
	}

	return feature, resp, nil
}

func (m *FeatureResource) ListFeatureDependencies(ctx context.Context, featureId string) ([]*Feature, *Response, error) {
	url := fmt.Sprintf("/api/v1/features/%v/dependencies", featureId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var feature []*Feature

	resp, err := rq.Do(ctx, req, &feature)
	if err != nil {
		return nil, resp, err
	}

	return feature, resp, nil
}

func (m *FeatureResource) ListFeatureDependents(ctx context.Context, featureId string) ([]*Feature, *Response, error) {
	url := fmt.Sprintf("/api/v1/features/%v/dependents", featureId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var feature []*Feature

	resp, err := rq.Do(ctx, req, &feature)
	if err != nil {
		return nil, resp, err
	}

	return feature, resp, nil
}

func (m *FeatureResource) UpdateFeatureLifecycle(ctx context.Context, featureId string, lifecycle string, qp *query.Params) (*Feature, *Response, error) {
	url := fmt.Sprintf("/api/v1/features/%v/%v", featureId, lifecycle)
	if qp != nil {
		url = url + qp.String()
	}

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var feature *Feature

	resp, err := rq.Do(ctx, req, &feature)
	if err != nil {
		return nil, resp, err
	}

	return feature, resp, nil
}
