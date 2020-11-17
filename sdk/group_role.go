package sdk

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type (
	Role struct {
		AssignmentType string `json:"assignmentType,omitempty"`
		Id             string `json:"id,omitempty"`
		Status         string `json:"status,omitempty"`
		Type           string `json:"type,omitempty"`
	}
)

var ValidAdminRoles = []string{"SUPER_ADMIN", "ORG_ADMIN", "API_ACCESS_MANAGEMENT_ADMIN", "APP_ADMIN", "USER_ADMIN", "MOBILE_ADMIN", "READ_ONLY_ADMIN", "HELP_DESK_ADMIN", "REPORT_ADMIN", "GROUP_MEMBERSHIP_ADMIN"}

func (m *ApiSupplement) DeleteAdminRole(id, roleID string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/groups/%s/roles/%s", id, roleID)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(context.Background(), req, nil)
}

func (m *ApiSupplement) ListAdminRoles(groupID string, qp *query.Params) (roles []*Role, resp *okta.Response, err error) {
	url := fmt.Sprintf("/api/v1/groups/%s/roles", groupID)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	resp, err = m.RequestExecutor.Do(context.Background(), req, &roles)
	return
}
func (m *ApiSupplement) CreateAdminRole(groupID string, body *Role, qp *query.Params) (*Role, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/groups/%s/roles", groupID)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	respBody := &Role{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, respBody)
	return respBody, resp, err
}
