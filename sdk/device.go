package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type Device struct {
	ID                  string                              `json:"id"`
	Status              string                              `json:"status"`
	Profile             map[string]interface{}              `json:"profile"`
	ResourceType        string                              `json:"resourceType"`
	ResourceDisplayName map[string]interface{}              `json:"sourceDisplayName"`
	ResourceAlternateID string                              `json:"resourceAlternateId"`
	ResourceID          string                              `json:"resourceId"`
	Links               interface{}                         `json:"_links"`
	Embedded            map[string][]map[string]interface{} `json:"_embedded"`
}

// ListDevices Gets all devices based on the query params
func (m *APISupplement) ListDevices(ctx context.Context, qp *query.Params) ([]*Device, *okta.Response, error) {
	url := "/api/v1/devices"
	if qp != nil {
		url += qp.String()
	}
    fmt.Printf("THE URL IS %+v", url)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var devices []*Device
	resp, err := m.RequestExecutor.Do(ctx, req, &devices)
	if err != nil {
		return nil, resp, err
	}
	return devices, resp, nil
}

// GetDevice gets device by ID
func (m *APISupplement) GetDevice(ctx context.Context, id string) (*Device, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/devices/%s", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var device *Device
	resp, err := m.RequestExecutor.Do(ctx, req, &device)
	if err != nil {
		return nil, resp, err
	}
	return device, resp, nil
}

// DeleteDevice deletes device by ID
func (m *APISupplement) DeleteDevice(ctx context.Context, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/devices/%s", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *APISupplement) ActivateDevice(ctx context.Context, id string) (*okta.Response, error) {
	return m.changeDeviceLifecycle(ctx, id, "activate")
}

func (m *APISupplement) DeactivateDevice(ctx context.Context, id string) (*okta.Response, error) {
	return m.changeDeviceLifecycle(ctx, id, "deactivate")
}

func (m *APISupplement) SuspendDevice(ctx context.Context, id string) (*okta.Response, error) {
	return m.changeDeviceLifecycle(ctx, id, "suspend")
}

func (m *APISupplement) UnsuspendDevice(ctx context.Context, id string) (*okta.Response, error) {
	return m.changeDeviceLifecycle(ctx, id, "unsuspend")
}

func (m *APISupplement) changeDeviceLifecycle(ctx context.Context, id, action string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/devices/%s/lifecycle/%s", id, action)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := m.RequestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
