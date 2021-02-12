package sdk

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

// Get role assigned to a user
func (m *ApiSupplement) GetUserAssignedRole(ctx context.Context, userId, roleId string) (*okta.Role, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/users/%v/roles/%v", userId, roleId)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	var role okta.Role
	resp, err := m.RequestExecutor.Do(ctx, req, &role)
	if err != nil {
		return nil, resp, err
	}
	return &role, resp, nil
}
