package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/okta/terraform-provider-okta/sdk/query"
)

type NetworkZoneResource resource

type NetworkZone struct {
	Links       interface{}            `json:"_links,omitempty"`
	Asns        []string               `json:"asns,omitempty"`
	Created     *time.Time             `json:"created,omitempty"`
	Gateways    []*NetworkZoneAddress  `json:"gateways,omitempty"`
	Id          string                 `json:"id,omitempty"`
	LastUpdated *time.Time             `json:"lastUpdated,omitempty"`
	Locations   []*NetworkZoneLocation `json:"locations,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Proxies     []*NetworkZoneAddress  `json:"proxies,omitempty"`
	ProxyType   string                 `json:"proxyType,omitempty"`
	Status      string                 `json:"status,omitempty"`
	System      *bool                  `json:"system,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Usage       string                 `json:"usage,omitempty"`
}

// Fetches a network zone from your Okta organization by &#x60;id&#x60;.
func (m *NetworkZoneResource) GetNetworkZone(ctx context.Context, zoneId string) (*NetworkZone, *Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%v", zoneId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var networkZone *NetworkZone

	resp, err := rq.Do(ctx, req, &networkZone)
	if err != nil {
		return nil, resp, err
	}

	return networkZone, resp, nil
}

// Updates a network zone in your organization.
func (m *NetworkZoneResource) UpdateNetworkZone(ctx context.Context, zoneId string, body NetworkZone) (*NetworkZone, *Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%v", zoneId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var networkZone *NetworkZone

	resp, err := rq.Do(ctx, req, &networkZone)
	if err != nil {
		return nil, resp, err
	}

	return networkZone, resp, nil
}

// Removes network zone.
func (m *NetworkZoneResource) DeleteNetworkZone(ctx context.Context, zoneId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%v", zoneId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Enumerates network zones added to your organization with pagination. A subset of zones can be returned that match a supported filter expression or query.
func (m *NetworkZoneResource) ListNetworkZones(ctx context.Context, qp *query.Params) ([]*NetworkZone, *Response, error) {
	url := "/api/v1/zones"
	if qp != nil {
		url = url + qp.String()
	}

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var networkZone []*NetworkZone

	resp, err := rq.Do(ctx, req, &networkZone)
	if err != nil {
		return nil, resp, err
	}

	return networkZone, resp, nil
}

// Adds a new network zone to your Okta organization.
func (m *NetworkZoneResource) CreateNetworkZone(ctx context.Context, body NetworkZone) (*NetworkZone, *Response, error) {
	url := "/api/v1/zones"

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var networkZone *NetworkZone

	resp, err := rq.Do(ctx, req, &networkZone)
	if err != nil {
		return nil, resp, err
	}

	return networkZone, resp, nil
}

// Activate Network Zone
func (m *NetworkZoneResource) ActivateNetworkZone(ctx context.Context, zoneId string) (*NetworkZone, *Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%v/lifecycle/activate", zoneId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var networkZone *NetworkZone

	resp, err := rq.Do(ctx, req, &networkZone)
	if err != nil {
		return nil, resp, err
	}

	return networkZone, resp, nil
}

// Deactivates a network zone.
func (m *NetworkZoneResource) DeactivateNetworkZone(ctx context.Context, zoneId string) (*NetworkZone, *Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%v/lifecycle/deactivate", zoneId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var networkZone *NetworkZone

	resp, err := rq.Do(ctx, req, &networkZone)
	if err != nil {
		return nil, resp, err
	}

	return networkZone, resp, nil
}
