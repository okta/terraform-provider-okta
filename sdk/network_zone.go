package sdk

import (
	"fmt"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
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
		Blacklist bool          `json:useAsBlackList,omitempty`
		Type      string        `json:"type,omitempty"`
	}
)

func (m *ApiSupplement) CreateNetworkZone(body NetworkZone, qp *query.Params) (*NetworkZone, *okta.Response, error) {
	url := "/api/v1/zones"
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	zone := body
	resp, err := m.RequestExecutor.Do(req, &zone)
	return &zone, resp, err
}

func (m *ApiSupplement) GetNetworkZone(id string) (*NetworkZone, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%s", id)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	zone := &NetworkZone{}
	resp, err := m.RequestExecutor.Do(req, &zone)
	if err != nil {
		return nil, resp, err
	}
	return zone, resp, nil
}

func (m *ApiSupplement) DeleteNetworkZone(id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%s", id)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(req, nil)
}

func (m *ApiSupplement) UpdateNetworkZone(id string, body NetworkZone, qp *query.Params) (*NetworkZone, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%s", id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	zone := body
	resp, err := m.RequestExecutor.Do(req, &zone)
	if err != nil {
		return nil, resp, err
	}
	return &zone, resp, nil
}
