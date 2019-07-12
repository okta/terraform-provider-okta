package okta

import (
	"fmt"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

type (
	Role struct {
		AssignmentType string `json:"assignmentType,omitempty"`
		Id             string `json:"id,omitempty"`
		Status         string `json:"status,omitempty"`
		Type           string `json:"type,omitempty"`
	}
)

var validAdminRoles = []string{"SUPER_ADMIN", "ORG_ADMIN", "API_ACCESS_MANAGEMENT_ADMIN", "APP_ADMIN", "USER_ADMIN", "MOBILE_ADMIN", "READ_ONLY_ADMIN", "HELP_DESK_ADMIN"}

func (m *ApiSupplement) DeleteAdminRole(id, roleId string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/groups/%s/roles/%s", id, roleId)
	req, err := m.requestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.requestExecutor.Do(req, nil)
}

func (m *ApiSupplement) ListAdminRoles(groupId string, qp *query.Params) (roles []*Role, resp *okta.Response, err error) {
	url := fmt.Sprintf("/api/v1/groups/%s/roles", groupId)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	resp, err = m.requestExecutor.Do(req, &roles)
	return
}
func (m *ApiSupplement) CreateAdminRole(groupId string, body *Role, qp *query.Params) (*Role, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/groups/%s/roles", groupId)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.requestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	respBody := &Role{}
	resp, err := m.requestExecutor.Do(req, respBody)
	return respBody, resp, err
}
