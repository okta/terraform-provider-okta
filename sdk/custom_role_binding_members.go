package sdk

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/okta/terraform-provider-okta/sdk/query"
)

type AddCustomRoleBindingMemberRequest struct {
	Additions []string `json:"additions,omitempty"`
}

type ListCustomRoleBindingMembersResponse struct {
	Members []*CustomRoleBindingMember `json:"members,omitempty"`
	Links   interface{}                `json:"_links,omitempty"`
}

type CustomRoleBindingMember struct {
	Id          string      `json:"id,omitempty"`
	Created     time.Time   `json:"created,omitempty"`
	LastUpdated time.Time   `json:"lastUpdated,omitempty"`
	Links       interface{} `json:"_links,omitempty"`
}

func (m *APISupplement) ListResourceSetBindingMembers(ctx context.Context, resourceSetID, customRoleID string, qp *query.Params) (*ListCustomRoleBindingMembersResponse, *Response, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s/bindings/%s/members", resourceSetID, customRoleID)
	if qp != nil {
		url += qp.String()
	}
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var members *ListCustomRoleBindingMembersResponse
	resp, err := re.Do(ctx, req, &members)
	if err != nil {
		return nil, resp, err
	}
	return members, resp, nil
}

func (m *APISupplement) AddResourceSetBindingMembers(ctx context.Context, resourceSetID, customRoleID string, body AddCustomRoleBindingMemberRequest) (*Response, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s/bindings/%s/members", resourceSetID, customRoleID)
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}

func (m *APISupplement) DeleteResourceSetBindingMember(ctx context.Context, resourceSetID, customRoleID, memberID string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s/bindings/%s/members/%s", resourceSetID, customRoleID, memberID)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}
