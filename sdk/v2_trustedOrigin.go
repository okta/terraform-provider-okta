package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/okta/terraform-provider-okta/sdk/query"
)

type TrustedOriginResource resource

type TrustedOrigin struct {
	Links         interface{} `json:"_links,omitempty"`
	Created       *time.Time  `json:"created,omitempty"`
	CreatedBy     string      `json:"createdBy,omitempty"`
	Id            string      `json:"id,omitempty"`
	LastUpdated   *time.Time  `json:"lastUpdated,omitempty"`
	LastUpdatedBy string      `json:"lastUpdatedBy,omitempty"`
	Name          string      `json:"name,omitempty"`
	Origin        string      `json:"origin,omitempty"`
	Scopes        []*Scope    `json:"scopes,omitempty"`
	Status        string      `json:"status,omitempty"`
}

func (m *TrustedOriginResource) CreateOrigin(ctx context.Context, body TrustedOrigin) (*TrustedOrigin, *Response, error) {
	url := "/api/v1/trustedOrigins"

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var trustedOrigin *TrustedOrigin

	resp, err := rq.Do(ctx, req, &trustedOrigin)
	if err != nil {
		return nil, resp, err
	}

	return trustedOrigin, resp, nil
}

func (m *TrustedOriginResource) GetOrigin(ctx context.Context, trustedOriginId string) (*TrustedOrigin, *Response, error) {
	url := fmt.Sprintf("/api/v1/trustedOrigins/%v", trustedOriginId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var trustedOrigin *TrustedOrigin

	resp, err := rq.Do(ctx, req, &trustedOrigin)
	if err != nil {
		return nil, resp, err
	}

	return trustedOrigin, resp, nil
}

func (m *TrustedOriginResource) UpdateOrigin(ctx context.Context, trustedOriginId string, body TrustedOrigin) (*TrustedOrigin, *Response, error) {
	url := fmt.Sprintf("/api/v1/trustedOrigins/%v", trustedOriginId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var trustedOrigin *TrustedOrigin

	resp, err := rq.Do(ctx, req, &trustedOrigin)
	if err != nil {
		return nil, resp, err
	}

	return trustedOrigin, resp, nil
}

func (m *TrustedOriginResource) DeleteOrigin(ctx context.Context, trustedOriginId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/trustedOrigins/%v", trustedOriginId)

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

func (m *TrustedOriginResource) ListOrigins(ctx context.Context, qp *query.Params) ([]*TrustedOrigin, *Response, error) {
	url := "/api/v1/trustedOrigins"
	if qp != nil {
		url = url + qp.String()
	}

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var trustedOrigin []*TrustedOrigin

	resp, err := rq.Do(ctx, req, &trustedOrigin)
	if err != nil {
		return nil, resp, err
	}

	return trustedOrigin, resp, nil
}

func (m *TrustedOriginResource) ActivateOrigin(ctx context.Context, trustedOriginId string) (*TrustedOrigin, *Response, error) {
	url := fmt.Sprintf("/api/v1/trustedOrigins/%v/lifecycle/activate", trustedOriginId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var trustedOrigin *TrustedOrigin

	resp, err := rq.Do(ctx, req, &trustedOrigin)
	if err != nil {
		return nil, resp, err
	}

	return trustedOrigin, resp, nil
}

func (m *TrustedOriginResource) DeactivateOrigin(ctx context.Context, trustedOriginId string) (*TrustedOrigin, *Response, error) {
	url := fmt.Sprintf("/api/v1/trustedOrigins/%v/lifecycle/deactivate", trustedOriginId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var trustedOrigin *TrustedOrigin

	resp, err := rq.Do(ctx, req, &trustedOrigin)
	if err != nil {
		return nil, resp, err
	}

	return trustedOrigin, resp, nil
}
