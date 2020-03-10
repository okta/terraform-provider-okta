package sdk

import (
	"encoding/json"
	"fmt"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
)

type (
	AuthScheme struct {
		Key   string `json:"key,omitempty"`
		Type  string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
	}

	Channel struct {
		Config  *HookConfig `json:"config"`
		Type    string      `json:"type,omitempty"`
		Version string      `json:"version,omitempty"`
	}

	HookConfig struct {
		AuthScheme *AuthScheme `json:"authScheme,omitempty"`
		Headers    []*Header   `json:"headers,omitempty"`
		URI        string      `json:"uri,omitempty"`
		Method     string      `json:"method,omitempty"`
	}

	Header struct {
		Key   string `json:"key,omitempty"`
		Value string `json:"value,omitempty"`
	}

	Events struct {
		Items []string `json:"items"`
	}

	EventHook struct {
		Channel *Channel `json:"channel"`
		ID      string   `json:"id,omitempty"`
		Name    string   `json:"name,omitempty"`
		Status  string   `json:"status,omitempty"`
		Events  *Events  `json:"events,omitempty"`
	}

	InlineHook struct {
		Channel *Channel `json:"channel"`
		ID      string   `json:"id,omitempty"`
		Name    string   `json:"name,omitempty"`
		Status  string   `json:"status,omitempty"`
		Type    string   `json:"type,omitempty"`
		Version string   `json:"version,omitempty"`
	}
)

func (e *Events) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type  string   `json:"type"`
		Items []string `json:"items"`
	}{
		Type:  "EVENT_TYPE",
		Items: e.Items,
	})
}

func activateHook(hookType string, id string, executor *okta.RequestExecutor) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/%sHooks/%s/lifecycle/activate", hookType, id)
	req, err := executor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	return executor.Do(req, nil)
}

func (m *ApiSupplement) ActivateInlineHook(id string) (*okta.Response, error) {
	return activateHook("inline", id, m.RequestExecutor)
}

func (m *ApiSupplement) ActivateEventHook(id string) (*okta.Response, error) {
	return activateHook("event", id, m.RequestExecutor)
}

func deactivateHook(hookType string, id string, executor *okta.RequestExecutor) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/%sHooks/%s/lifecycle/deactivate", hookType, id)
	req, err := executor.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	return executor.Do(req, nil)
}

func (m *ApiSupplement) DeactivateInlineHook(id string) (*okta.Response, error) {
	return deactivateHook("inline", id, m.RequestExecutor)
}

func (m *ApiSupplement) DeactivateEventHook(id string) (*okta.Response, error) {
	return deactivateHook("event", id, m.RequestExecutor)
}

func deleteHook(hookType string, id string, executor *okta.RequestExecutor) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/%sHooks/%s", hookType, id)
	req, err := executor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return executor.Do(req, nil)
}

func (m *ApiSupplement) DeleteInlineHook(id string) (*okta.Response, error) {
	return deleteHook("inline", id, m.RequestExecutor)
}

func (m *ApiSupplement) DeleteEventHook(id string) (*okta.Response, error) {
	return deleteHook("event", id, m.RequestExecutor)
}

func listHooks(hookType string, out interface{}, executor *okta.RequestExecutor) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/%sHooks", hookType)
	req, err := executor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := executor.Do(req, &out)
	return resp, err
}

func (m *ApiSupplement) ListInlineHooks() ([]*InlineHook, *okta.Response, error) {
	var hooks []*InlineHook
	resp, err := listHooks("inline", &hooks, m.RequestExecutor)
	return hooks, resp, err
}

func (m *ApiSupplement) ListEventHooks() ([]*EventHook, *okta.Response, error) {
	var hooks []*EventHook
	resp, err := listHooks("event", &hooks, m.RequestExecutor)
	return hooks, resp, err
}

func createHook(hookType string, body interface{}, qp *query.Params, executor *okta.RequestExecutor) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/%sHooks", hookType)
	if qp != nil {
		url += qp.String()
	}
	req, err := executor.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	resp, err := executor.Do(req, &body)
	return resp, err
}

func (m *ApiSupplement) CreateInlineHook(body InlineHook, qp *query.Params) (*InlineHook, *okta.Response, error) {
	resp, err := createHook("inline", &body, qp, m.RequestExecutor)
	return &body, resp, err
}

func (m *ApiSupplement) CreateEventHook(body EventHook, qp *query.Params) (*EventHook, *okta.Response, error) {
	resp, err := createHook("event", &body, qp, m.RequestExecutor)
	return &body, resp, err
}

func updateHook(hookType string, id string, body interface{}, qp *query.Params, executor *okta.RequestExecutor) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/%sHooks/%s", hookType, id)
	if qp != nil {
		url += qp.String()
	}
	req, err := executor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}

	return executor.Do(req, &body)
}

func (m *ApiSupplement) UpdateInlineHook(id string, body InlineHook, qp *query.Params) (*InlineHook, *okta.Response, error) {
	resp, err := updateHook("inline", id, &body, qp, m.RequestExecutor)
	return &body, resp, err
}

func (m *ApiSupplement) UpdateEventHook(id string, body EventHook, qp *query.Params) (*EventHook, *okta.Response, error) {
	resp, err := updateHook("event", id, &body, qp, m.RequestExecutor)
	return &body, resp, err
}

func getHook(hookType string, id string, hook interface{}, executor *okta.RequestExecutor) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/%sHooks/%s", hookType, id)
	req, err := executor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return executor.Do(req, hook)
}

func (m *ApiSupplement) GetInlineHook(id string) (*InlineHook, *okta.Response, error) {
	var hook InlineHook
	resp, err := getHook("inline", id, &hook, m.RequestExecutor)
	return &hook, resp, err
}

func (m *ApiSupplement) GetEventHook(id string) (*EventHook, *okta.Response, error) {
	var hook EventHook
	resp, err := getHook("event", id, &hook, m.RequestExecutor)
	return &hook, resp, err
}
