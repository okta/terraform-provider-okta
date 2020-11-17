package sdk

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type (
	EventHookAuthScheme struct {
		Key   string `json:"key,omitempty"`
		Type  string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
	}

	EventHookChannel struct {
		Config  *EventHookChannelConfig `json:"config"`
		Type    string                  `json:"type,omitempty"`
		Version string                  `json:"version,omitempty"`
	}

	EventHookChannelConfig struct {
		AuthScheme *EventHookAuthScheme `json:"authScheme,omitempty"`
		Headers    []*EventHookHeader   `json:"headers,omitempty"`
		URI        string               `json:"uri,omitempty"`
	}

	EventHookHeader struct {
		Key   string `json:"key,omitempty"`
		Value string `json:"value,omitempty"`
	}

	EventHookEvents struct {
		Type  string   `json:"type"`
		Items []string `json:"items"`
	}

	EventHook struct {
		Channel *EventHookChannel `json:"channel"`
		ID      string            `json:"id,omitempty"`
		Name    string            `json:"name,omitempty"`
		Status  string            `json:"status,omitempty"`
		Events  *EventHookEvents  `json:"events,omitempty"`
	}
)

func (m *ApiSupplement) ActivateEventHook(id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/eventHooks/%s/lifecycle/activate", id)
	req, err := m.RequestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(context.Background(), req, nil)
}

func (m *ApiSupplement) DeactivateEventHook(id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/eventHooks/%s/lifecycle/deactivate", id)
	req, err := m.RequestExecutor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(context.Background(), req, nil)
}

func (m *ApiSupplement) DeleteEventHook(id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/eventHooks/%s", id)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return m.RequestExecutor.Do(context.Background(), req, nil)
}

func (m *ApiSupplement) ListEventHooks() ([]*EventHook, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("GET", "/api/v1/eventHooks", nil)
	if err != nil {
		return nil, nil, err
	}

	var auth []*EventHook
	resp, err := m.RequestExecutor.Do(context.Background(), req, &auth)
	return auth, resp, err
}

func (m *ApiSupplement) CreateEventHook(body EventHook, qp *query.Params) (*EventHook, *okta.Response, error) {
	url := "/api/v1/eventHooks"
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	hook := body
	resp, err := m.RequestExecutor.Do(context.Background(), req, &hook)
	return &hook, resp, err
}

func (m *ApiSupplement) UpdateEventHook(id string, body EventHook, qp *query.Params) (*EventHook, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/eventHooks/%s", id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	hook := body
	resp, err := m.RequestExecutor.Do(context.Background(), req, &hook)
	if err != nil {
		return nil, resp, err
	}
	return &hook, resp, nil
}

func (m *ApiSupplement) GetEventHook(id string) (*EventHook, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/eventHooks/%s", id)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	hook := &EventHook{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, &hook)
	if err != nil {
		return nil, resp, err
	}
	return hook, resp, nil
}
