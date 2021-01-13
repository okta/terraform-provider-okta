package sdk

import (
	"context"
	"fmt"

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
		System    bool          `json:"system,omitempty"`
		Type      string        `json:"type,omitempty"`
		Usage     string        `json:"usage,omitempty"`
	}
)

func (m *ApiSupplement) CreateNetworkZone(ctx context.Context, body *NetworkZone, qp *query.Params) (*NetworkZone, *okta.Response, error) {
	url := "/api/v1/zones"
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	zone := body
	resp, err := m.RequestExecutor.Do(ctx, req, zone)
	return zone, resp, err
}

func (m *ApiSupplement) GetNetworkZone(ctx context.Context, id string) (*NetworkZone, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%s", id)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
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

func (m *ApiSupplement) ListNetworkZones(ctx context.Context) ([]*NetworkZone, *okta.Response, error) {
	url := "/api/v1/zones"
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
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

func (m *ApiSupplement) DeleteNetworkZone(ctx context.Context, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%s", id)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *ApiSupplement) UpdateNetworkZone(ctx context.Context, id string, body NetworkZone, qp *query.Params) (*NetworkZone, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%s", id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("PUT", url, body)
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
