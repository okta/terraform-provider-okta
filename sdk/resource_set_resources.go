// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/okta/terraform-provider-okta/sdk/query"
)

type AddResourceSetResourcesRequest struct {
	Additions []string `json:"additions"`
}

type ListResourceSetResourcesResponse struct {
	Resources []*ResourceSetResource `json:"resources"`
	Links     interface{}            `json:"_links"`
}

type ResourceSetResource struct {
	Id          string      `json:"id"`
	Created     time.Time   `json:"created"`
	LastUpdated time.Time   `json:"lastUpdated"`
	Links       interface{} `json:"_links"`
	Orn         string      `json:"orn"`
}

// ListResourceSetResources lists the resources that make up a Resource Set
func (m *APISupplement) ListResourceSetResources(ctx context.Context, resourceSetID string, qp *query.Params) (*ListResourceSetResourcesResponse, *Response, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s/resources", resourceSetID)
	if qp != nil {
		url += qp.String()
	}
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var resourceSets *ListResourceSetResourcesResponse
	resp, err := re.Do(ctx, req, &resourceSets)
	if err != nil {
		return nil, resp, err
	}
	return resourceSets, resp, nil
}

// AddResourceSetResources adds more resources to a Resource Set
func (m *APISupplement) AddResourceSetResources(ctx context.Context, resourceSetID string, body AddResourceSetResourcesRequest) (*Response, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s/resources", resourceSetID)
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}

// DeleteResourceSetResource Removes a resource from a Resource Set
func (m *APISupplement) DeleteResourceSetResource(ctx context.Context, resourceSetID, resourceID string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s/resources/%s", resourceSetID, resourceID)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}
