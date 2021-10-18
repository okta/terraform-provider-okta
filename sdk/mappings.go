package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type (
	MappingProperty struct {
		Expression string `json:"expression"`
		PushStatus string `json:"pushStatus"`
	}

	MappingSource struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}

	Mapping struct {
		ID         string                      `json:"id"`
		Source     *MappingSource              `json:"source,omitempty"`
		Target     *MappingSource              `json:"target,omitempty"`
		Properties map[string]*MappingProperty `json:"properties,omitempty"`
	}
)

func (m *APISupplement) GetProfileMappingBySourceId(ctx context.Context, sourceId, targetId string) (*Mapping, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/mappings?sourceId=%s&targetId=%s", sourceId, targetId)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var mappings []*Mapping
	resp, err := m.RequestExecutor.Do(ctx, req, &mappings)
	if err != nil {
		return nil, resp, err
	}
	for _, mapping := range mappings {
		if mapping.Source.ID == sourceId {
			return m.GetProfileMapping(ctx, mapping.ID)
		}
	}
	return nil, resp, err
}

func (m *APISupplement) GetProfileMapping(ctx context.Context, mappingId string) (*Mapping, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/mappings/%s", mappingId)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var mapping *Mapping
	resp, err := m.RequestExecutor.Do(ctx, req, &mapping)
	if err != nil {
		return nil, resp, err
	}
	return mapping, resp, err
}

func (m *APISupplement) UpdateMapping(ctx context.Context, mappingId string, body Mapping, qp *query.Params) (*Mapping, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/mappings/%s", mappingId)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	mapping := body
	resp, err := m.RequestExecutor.Do(ctx, req, &mapping)
	if err != nil {
		return nil, resp, err
	}
	return &mapping, resp, nil
}

// FindProfileMappingSource retrieves profile mapping source/target via name
func (m *APISupplement) FindProfileMappingSource(ctx context.Context, name, typ string, qp *query.Params) (*MappingSource, error) {
	uri := "/api/v1/mappings"
	if qp != nil {
		uri += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	var mappings []*Mapping
	resp, err := m.RequestExecutor.Do(ctx, req, &mappings)
	if err != nil {
		return nil, err
	}
	for {
		for _, m := range mappings {
			if m.Target.Name == name && m.Target.Type == typ {
				return m.Target, nil
			} else if m.Source.Name == name && m.Source.Type == typ {
				return m.Source, nil
			}
		}
		if resp.HasNextPage() {
			resp, err = resp.Next(ctx, &mappings)
			if err != nil {
				return nil, err
			}
			continue
		} else {
			break
		}
	}
	return nil, fmt.Errorf("could not locate profile mapping source with name '%s' and type '%s'", name, typ)
}
