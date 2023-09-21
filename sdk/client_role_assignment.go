package sdk

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Regular roles struct is missing `resource-set` and `role` fields, from CUSTOM role response.
type ClientRole struct {
	Embedded       interface{} `json:"_embedded,omitempty"`
	Links          interface{} `json:"_links,omitempty"`
	AssignmentType string      `json:"assignmentType,omitempty"`
	Created        *time.Time  `json:"created,omitempty"`
	Description    string      `json:"description,omitempty"`
	Id             string      `json:"id,omitempty"`
	Label          string      `json:"label,omitempty"`
	LastUpdated    *time.Time  `json:"lastUpdated,omitempty"`
	Status         string      `json:"status,omitempty"`
	Type           string      `json:"type,omitempty"`
	ResourceSet    string      `json:"resource-set,omitempty"`
	Role           string      `json:"role,omitempty"`
}

func (m *APISupplement) ListClientRoles(ctx context.Context, clientID string) ([]*ClientRole, *Response, error) {
	var roles []*ClientRole

	url := fmt.Sprintf("/oauth2/v1/clients/%s/roles", clientID)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := m.RequestExecutor.Do(ctx, req, &roles)
	return roles, resp, err
}

type ClientRoleAssignment struct {
	ResourceSet string `json:"resource-set,omitempty"`
	Role        string `json:"role,omitempty"`
	Type        string `json:"type"`
}

func (m *APISupplement) AssignClientRole(ctx context.Context, clientID string, assignment *ClientRoleAssignment) (*ClientRole, *Response, error) {
	var role *ClientRole

	url := fmt.Sprintf("/oauth2/v1/clients/%s/roles", clientID)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, assignment)
	if err != nil {
		return nil, nil, err
	}
	resp, err := m.RequestExecutor.Do(ctx, req, &role)
	return role, resp, err
}

func (m *APISupplement) GetClientRole(ctx context.Context, clientID, roleID string) (*ClientRole, *Response, error) {
	var role *ClientRole

	url := fmt.Sprintf("/oauth2/v1/clients/%s/roles/%s", clientID, roleID)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := m.RequestExecutor.Do(ctx, req, &role)
	return role, resp, err
}

func (m *APISupplement) UnassignClientRole(ctx context.Context, clientID, roleID string) (*Response, error) {
	url := fmt.Sprintf("/oauth2/v1/clients/%s/roles/%s", clientID, roleID)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := m.RequestExecutor.Do(ctx, req, nil)
	return resp, err
}
