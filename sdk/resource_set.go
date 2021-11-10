package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type ResourceSet struct {
	Id          string   `json:"id,omitempty"`
	Label       string   `json:"label,omitempty"`
	Description string   `json:"description,omitempty"`
	Resources   []string `json:"resources,omitempty"`
}

type ListResourceSetsResponse struct {
	ResourceSets []*ResourceSet `json:"resource-sets,omitempty"`
}

// ListResourceSets Gets all ResourceSets
func (m *APISupplement) ListResourceSets(ctx context.Context) (*ListResourceSetsResponse, *okta.Response, error) {
	url := "/api/v1/iam/resource-sets"
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var resourceSets *ListResourceSetsResponse
	resp, err := re.Do(ctx, req, &resourceSets)
	if err != nil {
		return nil, resp, err
	}
	return resourceSets, resp, nil
}

// GetResourceSet gets ResourceSet by ID
func (m *APISupplement) GetResourceSet(ctx context.Context, resourceSetID string) (*ResourceSet, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s", resourceSetID)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var resourceSet *ResourceSet
	resp, err := re.Do(ctx, req, &resourceSet)
	if err != nil {
		return nil, resp, err
	}
	return resourceSet, resp, nil
}

// CreateResourceSet creates ResourceSet
func (m *APISupplement) CreateResourceSet(ctx context.Context, body ResourceSet) (*ResourceSet, *okta.Response, error) {
	url := "/api/v1/iam/resource-sets"
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var resourceSet *ResourceSet
	resp, err := re.Do(ctx, req, &resourceSet)
	if err != nil {
		return nil, resp, err
	}
	return resourceSet, resp, nil
}

// UpdateResourceSet updates ResourceSet
func (m *APISupplement) UpdateResourceSet(ctx context.Context, resourceSetID string, body ResourceSet) (*ResourceSet, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s", resourceSetID)
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}
	var resourceSet *ResourceSet
	resp, err := re.Do(ctx, req, &resourceSet)
	if err != nil {
		return nil, resp, err
	}
	return resourceSet, resp, nil
}

// DeleteResourceSet deletes ResourceSet by ID
func (m *APISupplement) DeleteResourceSet(ctx context.Context, resourceSetID string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s", resourceSetID)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}
