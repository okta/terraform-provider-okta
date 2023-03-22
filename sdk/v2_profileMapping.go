package sdk

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type ProfileMappingResource resource

type ProfileMapping struct {
	Links      interface{}                        `json:"_links,omitempty"`
	Id         string                             `json:"id,omitempty"`
	Properties map[string]*ProfileMappingProperty `json:"properties,omitempty"`
	Source     *ProfileMappingSource              `json:"source,omitempty"`
	Target     *ProfileMappingSource              `json:"target,omitempty"`
}

// Fetches a single Profile Mapping referenced by its ID.
func (m *ProfileMappingResource) GetProfileMapping(ctx context.Context, mappingId string) (*ProfileMapping, *Response, error) {
	url := fmt.Sprintf("/api/v1/mappings/%v", mappingId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var profileMapping *ProfileMapping

	resp, err := rq.Do(ctx, req, &profileMapping)
	if err != nil {
		return nil, resp, err
	}

	return profileMapping, resp, nil
}

// Updates an existing Profile Mapping by adding, updating, or removing one or many Property Mappings.
func (m *ProfileMappingResource) UpdateProfileMapping(ctx context.Context, mappingId string, body ProfileMapping) (*ProfileMapping, *Response, error) {
	url := fmt.Sprintf("/api/v1/mappings/%v", mappingId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var profileMapping *ProfileMapping

	resp, err := rq.Do(ctx, req, &profileMapping)
	if err != nil {
		return nil, resp, err
	}

	return profileMapping, resp, nil
}

// Enumerates Profile Mappings in your organization with pagination.
func (m *ProfileMappingResource) ListProfileMappings(ctx context.Context, qp *query.Params) ([]*ProfileMapping, *Response, error) {
	url := fmt.Sprintf("/api/v1/mappings")
	if qp != nil {
		url = url + qp.String()
	}

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var profileMapping []*ProfileMapping

	resp, err := rq.Do(ctx, req, &profileMapping)
	if err != nil {
		return nil, resp, err
	}

	return profileMapping, resp, nil
}
