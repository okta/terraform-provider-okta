// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"fmt"
	"net/http"
)

type ResourceSet struct {
	Id          string   `json:"id,omitempty"`
	Label       string   `json:"label,omitempty"`
	Description string   `json:"description,omitempty"`
	Resources   []string `json:"resources,omitempty"`
}

// PatchResourceSet
// https://developer.okta.com/docs/api/openapi/okta-management/management/tag/ResourceSet/#tag/ResourceSet/operation/addResourceSetResource
type PatchResourceSet struct {
	Additions []string `json:"additions,omitempty"`
}

type ListResourceSetsResponse struct {
	ResourceSets []*ResourceSet `json:"resource-sets,omitempty"`
}

// ListResourceSets Gets all ResourceSets
func (m *APISupplement) ListResourceSets(ctx context.Context) (*ListResourceSetsResponse, *Response, error) {
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
func (m *APISupplement) GetResourceSet(ctx context.Context, resourceSetID string) (*ResourceSet, *Response, error) {
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
func (m *APISupplement) CreateResourceSet(ctx context.Context, body ResourceSet) (*ResourceSet, *Response, error) {
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
func (m *APISupplement) UpdateResourceSet(ctx context.Context, resourceSetID string, body ResourceSet) (*ResourceSet, *Response, error) {
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

// PatchResourceSet patches a ResourceSet with additional additions
// https://developer.okta.com/docs/api/openapi/okta-management/management/tag/ResourceSet/#tag/ResourceSet/operation/addResourceSetResource
func (m *APISupplement) PatchResourceSet(ctx context.Context, resourceSetID string, body PatchResourceSet) (*ResourceSet, *Response, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s/resources", resourceSetID)
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPatch, url, body)
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
func (m *APISupplement) DeleteResourceSet(ctx context.Context, resourceSetID string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s", resourceSetID)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}
