// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"fmt"
	"net/http"
)

type AppUserType struct {
	Id          string      `json:"id"`
	DisplayName string      `json:"displayName"`
	Name        string      `json:"name"`
	Schemas     []string    `json:"schemas"`
	IsDefault   bool        `json:"isDefault"`
	Type        string      `json:"type"`
	Links       interface{} `json:"_links"`
}

func (m *APISupplement) GetAppUserTypes(ctx context.Context, appID string) ([]*AppUserType, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%s/user/types", appID)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var appUserTypes []*AppUserType
	resp, err := m.RequestExecutor.Do(ctx, req, &appUserTypes)
	if err != nil {
		return nil, resp, err
	}
	return appUserTypes, resp, nil
}
