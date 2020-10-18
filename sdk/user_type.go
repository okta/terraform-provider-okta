package sdk

import (
	"fmt"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
	"net/url"
	"strings"
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

func (m *ApiSupplement) DeleteUserType(id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/meta/types/user/%s", id)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)

	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(req, nil)
}

func (m *ApiSupplement) ListUserTypes() ([]*UserType, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("GET", "/api/v1/meta/types/user/", nil)
	if err != nil {
		return nil, nil, err
	}

	var auth []*UserType
	resp, err := m.RequestExecutor.Do(req, &auth)
	return auth, resp, err
}

func (m *ApiSupplement) CreateUserType(body UserType, qp *query.Params) (*UserType, *okta.Response, error) {
	url := "/api/v1/meta/types/user/"
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	userType := body
	resp, err := m.RequestExecutor.Do(req, &userType)
	return &userType, resp, err
}

func (m *ApiSupplement) UpdateUserType(id string, body UserType, qp *query.Params) (*UserType, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/meta/types/user/%s", id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	userType := body
	resp, err := m.RequestExecutor.Do(req, &userType)
	if err != nil {
		return nil, resp, err
	}
	return &userType, resp, nil
}

func (m *ApiSupplement) GetUserType(id string) (*UserType, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/meta/types/user/%s", id)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	userType := &UserType{}
	resp, err := m.RequestExecutor.Do(req, &userType)
	if err != nil {
		return nil, resp, err
	}
	return userType, resp, nil
}

func (c *ApiSupplement) FindUserType(name string, qp *query.Params) (*UserType, error) {
	userTypeList, res, err := c.ListUserTypes()
	if err != nil {
		return nil, err
	}

	for _, userType := range userTypeList {
		if strings.EqualFold(name, userType.Name) {
			return userType, nil
		}
	}

	if after := getNextLinkOffset(res); after != "" {
		qp.After = after
		return c.FindUserType(name, qp)
	}
	return nil, nil
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
