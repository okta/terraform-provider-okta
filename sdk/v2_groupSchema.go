package sdk

import (
	"context"
)

type GroupSchemaResource resource

type GroupSchema struct {
	Schema      string                  `json:"$schema,omitempty"`
	Links       interface{}             `json:"_links,omitempty"`
	Created     string                  `json:"created,omitempty"`
	Definitions *GroupSchemaDefinitions `json:"definitions,omitempty"`
	Description string                  `json:"description,omitempty"`
	Id          string                  `json:"id,omitempty"`
	LastUpdated string                  `json:"lastUpdated,omitempty"`
	Name        string                  `json:"name,omitempty"`
	Properties  *UserSchemaProperties   `json:"properties,omitempty"`
	Title       string                  `json:"title,omitempty"`
	Type        string                  `json:"type,omitempty"`
}

// Fetches the group schema
func (m *GroupSchemaResource) GetGroupSchema(ctx context.Context) (*GroupSchema, *Response, error) {
	url := "/api/v1/meta/schemas/group/default"

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var groupSchema *GroupSchema

	resp, err := rq.Do(ctx, req, &groupSchema)
	if err != nil {
		return nil, resp, err
	}

	return groupSchema, resp, nil
}

// Updates, adds ore removes one or more custom Group Profile properties in the schema
func (m *GroupSchemaResource) UpdateGroupSchema(ctx context.Context, body GroupSchema) (*GroupSchema, *Response, error) {
	url := "/api/v1/meta/schemas/group/default"

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var groupSchema *GroupSchema

	resp, err := rq.Do(ctx, req, &groupSchema)
	if err != nil {
		return nil, resp, err
	}

	return groupSchema, resp, nil
}
