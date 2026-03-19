// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/terraform-provider-okta/sdk/query"
)

type ListCustomRolesResponse struct {
	Roles []*CustomRole `json:"roles,omitempty"`
}

type CustomRole struct {
	Id          string      `json:"id,omitempty"`
	Label       string      `json:"label,omitempty"`
	Description string      `json:"description,omitempty"`
	Permissions []string    `json:"permissions,omitempty"`
	Links       interface{} `json:"_links,omitempty"`
}

// ListCustomRoles Gets all customRoles based on the query params
func (m *APISupplement) ListCustomRoles(ctx context.Context, qp *query.Params) (*ListCustomRolesResponse, *Response, error) {
	url := "/api/v1/iam/roles"
	if qp != nil {
		url += qp.String()
	}
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var response *ListCustomRolesResponse
	resp, err := re.Do(ctx, req, &response)
	if err != nil {
		return nil, resp, err
	}
	return response, resp, nil
}

// GetCustomRole gets customRole by ID
func (m *APISupplement) GetCustomRole(ctx context.Context, roleIdOrLabel string) (*CustomRole, *Response, error) {
	url := fmt.Sprintf("/api/v1/iam/roles/%s", roleIdOrLabel)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var customRole *CustomRole
	resp, err := re.Do(ctx, req, &customRole)
	if err != nil {
		return nil, resp, err
	}
	return customRole, resp, nil
}

// CreateCustomRole creates customRole
func (m *APISupplement) CreateCustomRole(ctx context.Context, body CustomRole) (*CustomRole, *Response, error) {
	url := "/api/v1/iam/roles"
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var customRole *CustomRole
	resp, err := re.Do(ctx, req, &customRole)
	if err != nil {
		return nil, resp, err
	}
	return customRole, resp, nil
}

// UpdateCustomRole updates customRole
func (m *APISupplement) UpdateCustomRole(ctx context.Context, roleIdOrLabel string, body CustomRole) (*CustomRole, *Response, error) {
	url := fmt.Sprintf("/api/v1/iam/roles/%s", roleIdOrLabel)
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}
	var customRole *CustomRole
	resp, err := re.Do(ctx, req, &customRole)
	if err != nil {
		return nil, resp, err
	}
	return customRole, resp, nil
}

// DeleteCustomRole deletes customRole by ID
func (m *APISupplement) DeleteCustomRole(ctx context.Context, roleIdOrLabel string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/iam/roles/%s", roleIdOrLabel)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}
