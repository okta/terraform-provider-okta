package okta

import (
	"fmt"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

type (
	Gateways struct {
		Type  string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
	}

	NetworkZone struct {
		Gateways []*Gateways `json:"gateways"`
		ID       string      `json:"id,omitempty"`
		Name     string      `json:"name,omitempty"`
		System   bool        `json:"system,omitempty"`
		Type     string      `json:"type,omitempty"`
	}
)

func (m *ApiSupplement) CreateNetworkZone(body NetworkZone, qp *query.Params) (*NetworkZone, *okta.Response, error) {
	url := "/api/v1/zones"
	if qp != nil {
		url += qp.String()
	}
	req, err := m.requestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	zone := body
	resp, err := m.requestExecutor.Do(req, &zone)
	return &zone, resp, err
}

func (m *ApiSupplement) GetNetworkZone(id string) (*NetworkZone, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%s", id)
	req, err := m.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	zone := &NetworkZone{}
	resp, err := m.requestExecutor.Do(req, &zone)
	if err != nil {
		return nil, resp, err
	}
	return zone, resp, nil
}

func (m *ApiSupplement) DeleteNetworkZone(id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%s", id)
	req, err := m.requestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.requestExecutor.Do(req, nil)
}

func (m *ApiSupplement) UpdateNetworkZone(id string, body NetworkZone, qp *query.Params) (*NetworkZone, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/zones/%s", id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.requestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	zone := body
	resp, err := m.requestExecutor.Do(req, &zone)
	if err != nil {
		return nil, resp, err
	}
	return &zone, resp, nil
}
