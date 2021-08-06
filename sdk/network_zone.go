package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type (
	AddressObj struct {
		Type  string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
	}

	Location struct {
		Country string `json:"country,omitempty"`
		Region  string `json:"region,omitempty"`
	}

	NetworkZone struct {
		Gateways  []*AddressObj `json:"gateways,omitempty"`
		ID        string        `json:"id,omitempty"`
		Locations []*Location   `json:"locations,omitempty"`
		Name      string        `json:"name,omitempty"`
		Proxies   []*AddressObj `json:"proxies,omitempty"`
		ProxyType string        `json:"proxyType,omitempty"`
		System    bool          `json:"system,omitempty"`
		Type      string        `json:"type,omitempty"`
		Usage     string        `json:"usage,omitempty"`
	}
)

func (m *APISupplement) CreateNetworkZone(ctx context.Context, body *NetworkZone, qp *query.Params) (*NetworkZone, *okta.Response, error) {
	url := "/api/v1/zones"
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}

	zone := body
	resp, err := m.RequestExecutor.Do(ctx, req, zone)
	return zone, resp, err
}

func (m *APISupplement) GetNetworkZone(ctx context.Context, id string) (*NetworkZone, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%s", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	zone := &NetworkZone{}
	resp, err := m.RequestExecutor.Do(ctx, req, &zone)
	if err != nil {
		return nil, resp, err
	}
	return zone, resp, nil
}

func (m *APISupplement) ListNetworkZones(ctx context.Context) ([]*NetworkZone, *okta.Response, error) {
	url := "/api/v1/zones"
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var zones []*NetworkZone
	resp, err := m.RequestExecutor.Do(ctx, req, &zones)
	if err != nil {
		return nil, resp, err
	}
	return zones, resp, nil
}

func (m *APISupplement) DeleteNetworkZone(ctx context.Context, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%s", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *APISupplement) UpdateNetworkZone(ctx context.Context, id string, body NetworkZone, qp *query.Params) (*NetworkZone, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%s", id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}

	zone := body
	resp, err := m.RequestExecutor.Do(ctx, req, &zone)
	if err != nil {
		return nil, resp, err
	}
	return &zone, resp, nil
}
