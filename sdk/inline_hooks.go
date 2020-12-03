package sdk

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type (
	InlineHookAuthScheme struct {
		Key   string `json:"key,omitempty"`
		Type  string `json:"type,omitempty"`
		Value string `json:"value,omitempty"`
	}

	InlineHookChannel struct {
		Config  *InlineHookChannelConfig `json:"config"`
		Type    string                   `json:"type,omitempty"`
		Version string                   `json:"version,omitempty"`
	}

	InlineHookChannelConfig struct {
		AuthScheme *InlineHookAuthScheme `json:"authScheme,omitempty"`
		Headers    []*InlineHookHeader   `json:"headers,omitempty"`
		URI        string                `json:"uri,omitempty"`
		Method     string                `json:"method,omitempty"`
	}

	InlineHookHeader struct {
		Key   string `json:"key,omitempty"`
		Value string `json:"value,omitempty"`
	}

	InlineHook struct {
		Channel *InlineHookChannel `json:"channel"`
		ID      string             `json:"id,omitempty"`
		Name    string             `json:"name,omitempty"`
		Status  string             `json:"status,omitempty"`
		Type    string             `json:"type,omitempty"`
		Version string             `json:"version,omitempty"`
	}
)

func (m *ApiSupplement) CreateInlineHook(ctx context.Context, body InlineHook, qp *query.Params) (*InlineHook, *okta.Response, error) {
	url := "/api/v1/inlineHooks"
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	hook := body
	resp, err := m.RequestExecutor.Do(ctx, req, &hook)
	return &hook, resp, err
}

func (m *ApiSupplement) UpdateInlineHook(ctx context.Context, id string, body InlineHook, qp *query.Params) (*InlineHook, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/inlineHooks/%s", id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	hook := body
	resp, err := m.RequestExecutor.Do(ctx, req, &hook)
	if err != nil {
		return nil, resp, err
	}
	return &hook, resp, nil
}

func (m *ApiSupplement) GetInlineHook(ctx context.Context, id string) (*InlineHook, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/inlineHooks/%s", id)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	hook := &InlineHook{}
	resp, err := m.RequestExecutor.Do(ctx, req, &hook)
	if err != nil {
		return nil, resp, err
	}
	return hook, resp, nil
}

func (m *ApiSupplement) ListInlineHooks(ctx context.Context) ([]*InlineHook, *okta.Response, error) {
	req, err := m.RequestExecutor.NewRequest("GET", "/api/v1/inlineHooks", nil)
	if err != nil {
		return nil, nil, err
	}

	var auth []*InlineHook
	resp, err := m.RequestExecutor.Do(ctx, req, &auth)
	return auth, resp, err
}
