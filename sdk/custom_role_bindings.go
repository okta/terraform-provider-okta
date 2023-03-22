package sdk

import (
	"context"
	"fmt"
	"net/http"
)

type CreateCustomRoleBindingRequest struct {
	Role    string   `json:"role,omitempty"`
	Members []string `json:"members,omitempty"`
}

type CustomRoleBinding struct {
	Id    string      `json:"id,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}

func (m *APISupplement) CreateResourceSetBinding(ctx context.Context, resourceSetID string, body CreateCustomRoleBindingRequest) (*Response, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s/bindings", resourceSetID)
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}

func (m *APISupplement) GetResourceSetBinding(ctx context.Context, resourceSetID, customRoleID string) (*CustomRoleBinding, *Response, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s/bindings/%s", resourceSetID, customRoleID)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var customRoleBinding *CustomRoleBinding
	resp, err := re.Do(ctx, req, &customRoleBinding)
	if err != nil {
		return nil, resp, err
	}
	return customRoleBinding, resp, nil
}

func (m *APISupplement) DeleteResourceSetBinding(ctx context.Context, resourceSetID, customRoleID string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s/bindings/%s", resourceSetID, customRoleID)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}
