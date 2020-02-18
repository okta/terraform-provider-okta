package sdk

import (
	"fmt"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
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

func (m *ApiSupplement) RemovePropertyMapping(mappingId, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/mappings/%s/", mappingId)
	req, err := m.RequestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(req, nil)
}

func (m *ApiSupplement) GetProfileMappingBySourceId(sourceId, targetId string) (*Mapping, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/mappings?sourceId=%s&targetId=%s", sourceId, targetId)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)

	if err != nil {
		return nil, nil, err
	}

	var mappings []*Mapping
	resp, err := m.RequestExecutor.Do(req, &mappings)

	for _, mapping := range mappings {
		if mapping.Source.ID == sourceId {
			return m.GetProfileMapping(mapping.ID)
		}
	}

	return nil, resp, err
}

func (m *ApiSupplement) GetProfileMapping(mappingId string) (*Mapping, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/mappings/%s", mappingId)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	mapping := &Mapping{}
	resp, err := m.RequestExecutor.Do(req, mapping)
	return mapping, resp, err
}

func (m *ApiSupplement) AddPropertyMapping(mappingId string, body Mapping, qp *query.Params) (*Mapping, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/mappings/%s", mappingId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	mapping := body
	resp, err := m.RequestExecutor.Do(req, &mapping)
	return &mapping, resp, err
}

func (m *ApiSupplement) UpdateMapping(mappingId string, body Mapping, qp *query.Params) (*Mapping, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/mappings/%s", mappingId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	mapping := body
	resp, err := m.RequestExecutor.Do(req, &mapping)
	if err != nil {
		return nil, resp, err
	}
	return &mapping, resp, nil
}

// FindProfileMappingSource retrieves profile mapping source/target via name
func (m *ApiSupplement) FindProfileMappingSource(name, typ string, qp *query.Params) (*MappingSource, *okta.Response, error) {
	uri := "/api/v1/mappings"

	if qp != nil {
		uri += qp.String()
	}

	req, err := m.RequestExecutor.NewRequest("GET", uri, nil)

	if err != nil {
		return nil, nil, err
	}

	mappings := []*Mapping{}
	res, err := m.RequestExecutor.Do(req, &mappings)

	for _, m := range mappings {
		if m.Target.Name == name && m.Target.Type == typ {
			return m.Target, res, nil
		} else if m.Source.Name == name && m.Source.Type == typ {
			return m.Source, res, nil
		}
	}

	if after := GetAfterParam(res); after != "" {
		qp.After = after

		return m.FindProfileMappingSource(name, typ, qp)
	}

	return nil, res, fmt.Errorf("could not locate profile mapping source with name %s", name)
}
