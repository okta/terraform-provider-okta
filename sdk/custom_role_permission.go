package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type ListPermissionsResponse struct {
	Permissions []*Permission `json:"permissions,omitempty"`
}

type Permission struct {
	Label string      `json:"label"`
	Links interface{} `json:"links"`
}

func (m *APISupplement) ListCustomRolePermissions(ctx context.Context, roleIdOrLabel string) (*ListPermissionsResponse, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/iam/roles/%s/permissions", roleIdOrLabel)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var listPermissionsResponse *ListPermissionsResponse
	resp, err := re.Do(ctx, req, &listPermissionsResponse)
	if err != nil {
		return nil, resp, err
	}
	return listPermissionsResponse, resp, nil
}

func (m *APISupplement) AddCustomRolePermission(ctx context.Context, roleIdOrLabel, permissionType string) (*Permission, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/iam/roles/%s/permissions/%s", roleIdOrLabel, permissionType)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var permission *Permission
	resp, err := re.Do(ctx, req, &permission)
	if err != nil {
		return nil, resp, err
	}
	return permission, resp, nil
}

func (m *APISupplement) DeleteCustomRolePermission(ctx context.Context, roleIdOrLabel, permissionType string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/iam/roles/%s/permissions/%s", roleIdOrLabel, permissionType)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}
