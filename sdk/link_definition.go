package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type LinkedObject struct {
	Primary    *LinkedObjectPart `json:"primary"`
	Associated *LinkedObjectPart `json:"associated"`
}

type LinkedObjectPart struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

func (m *APISupplement) ListLinkedObjects(ctx context.Context) ([]*LinkedObject, *okta.Response, error) {
	url := "/api/v1/meta/schemas/user/linkedObjects"
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var linkedObjects []*LinkedObject
	resp, err := re.Do(ctx, req, &linkedObjects)
	if err != nil {
		return nil, resp, err
	}
	return linkedObjects, resp, nil
}

func (m *APISupplement) CreateLinkedObject(ctx context.Context, body LinkedObject) (*LinkedObject, *okta.Response, error) {
	url := "/api/v1/meta/schemas/user/linkedObjects"
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var linkedObject *LinkedObject
	resp, err := re.Do(ctx, req, &linkedObject)
	if err != nil {
		return nil, resp, err
	}
	return linkedObject, resp, nil
}

func (m *APISupplement) GetLinkedObject(ctx context.Context, name string) (*LinkedObject, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/meta/schemas/user/linkedObjects/%s", name)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var linkedObject *LinkedObject
	resp, err := re.Do(ctx, req, &linkedObject)
	if err != nil {
		return nil, resp, err
	}
	return linkedObject, resp, nil
}

func (m *APISupplement) DeleteLinkedObject(ctx context.Context, name string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/meta/schemas/user/linkedObjects/%s", name)
	re := m.cloneRequestExecutor()
	req, err := re.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}
