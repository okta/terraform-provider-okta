package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type LinkedObjectValues struct {
	Links interface{} `json:"_links,omitempty"`
}

func (m *APISupplement) SetLinkedObjectValueForPrimary(ctx context.Context, associatedUserId, primaryName, primaryUserId string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/users/%s/linkedObjects/%s/%s", associatedUserId, primaryName, primaryUserId)
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}

func (m *APISupplement) GetLinkedObjectValues(ctx context.Context, userId, primaryOrAssociatedName string) ([]*LinkedObjectValues, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/users/%s/linkedObjects/%s", userId, primaryOrAssociatedName)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var linkedObjectValues []*LinkedObjectValues
	resp, err := re.Do(ctx, req, &linkedObjectValues)
	if err != nil {
		return nil, resp, err
	}
	return linkedObjectValues, resp, nil
}

func (m *APISupplement) DeleteLinkedObjectValue(ctx context.Context, userId, primaryName string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/users/%s/linkedObjects/%s", userId, primaryName)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}
