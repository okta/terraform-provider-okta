package sdk

import (
	"context"
	"fmt"
)

type LinkedObjectResource resource

type LinkedObject struct {
	Links      interface{}          `json:"_links,omitempty"`
	Associated *LinkedObjectDetails `json:"associated,omitempty"`
	Primary    *LinkedObjectDetails `json:"primary,omitempty"`
}

func (m *LinkedObjectResource) AddLinkedObjectDefinition(ctx context.Context, body LinkedObject) (*LinkedObject, *Response, error) {
	url := "/api/v1/meta/schemas/user/linkedObjects"

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var linkedObject *LinkedObject

	resp, err := rq.Do(ctx, req, &linkedObject)
	if err != nil {
		return nil, resp, err
	}

	return linkedObject, resp, nil
}

func (m *LinkedObjectResource) GetLinkedObjectDefinition(ctx context.Context, linkedObjectName string) (*LinkedObject, *Response, error) {
	url := fmt.Sprintf("/api/v1/meta/schemas/user/linkedObjects/%v", linkedObjectName)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var linkedObject *LinkedObject

	resp, err := rq.Do(ctx, req, &linkedObject)
	if err != nil {
		return nil, resp, err
	}

	return linkedObject, resp, nil
}

func (m *LinkedObjectResource) DeleteLinkedObjectDefinition(ctx context.Context, linkedObjectName string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/meta/schemas/user/linkedObjects/%v", linkedObjectName)

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

func (m *LinkedObjectResource) ListLinkedObjectDefinitions(ctx context.Context) ([]*LinkedObject, *Response, error) {
	url := "/api/v1/meta/schemas/user/linkedObjects"

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var linkedObject []*LinkedObject

	resp, err := rq.Do(ctx, req, &linkedObject)
	if err != nil {
		return nil, resp, err
	}

	return linkedObject, resp, nil
}
