package sdk

import (
	"context"
	"fmt"
	"net/url"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type (
	UserTypeLinks struct {
		Schema struct {
			Href string `json:"href"`
		} `json:"schema"`
	}

	UserType struct {
		Id          string         `json:"id,omitempty"`
		Name        string         `json:"name,omitempty"`
		DisplayName string         `json:"displayName,omitempty"`
		Description string         `json:"description,omitempty"`
		Links       *UserTypeLinks `json:"_links,omitempty"`
	}
)

func (m *ApiSupplement) ListUserTypes() ([]*UserType, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("GET", "/api/v1/meta/types/user", nil)
	if err != nil {
		return nil, nil, err
	}

	var auth []*UserType
	resp, err := m.RequestExecutor.Do(context.Background(), req, &auth)
	return auth, resp, err
}

func (m *ApiSupplement) GetUserType(id string) (*UserType, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("GET", fmt.Sprintf("/api/v1/meta/types/user/%s", id), nil)
	if err != nil {
		return nil, nil, err
	}

	userType := &UserType{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, &userType)
	if err != nil {
		return nil, resp, err
	}
	return userType, resp, nil
}

func (c *ApiSupplement) GetUserTypeSchemaUrl(id string, qp *query.Params) (string, error) {
	if id == "" {
		id = "default"
	}
	userType, _, err := c.GetUserType(id)
	if err != nil {
		return "", err
	}

	if userType != nil {
		u, _ := url.Parse(userType.Links.Schema.Href)
		var href = u.EscapedPath()
		return href, nil
	}
	return "", nil
}
