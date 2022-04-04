package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type (
	EmailTemplate struct {
		Id              string                       `json:"id,omitempty"`
		Name            string                       `json:"name,omitempty"`
		Type            string                       `json:"type,omitempty"`
		DefaultLanguage string                       `json:"defaultLanguage,omitempty"`
		Subject         string                       `json:"subject,omitempty"`
		Template        string                       `json:"template,omitempty"`
		Translations    map[string]*EmailTranslation `json:"translations,omitempty"`
	}

	EmailTranslation struct {
		Subject  string `json:"subject,omitempty"`
		Template string `json:"template,omitempty"`
	}
)

func (m *APISupplement) CreateEmailTemplate(ctx context.Context, body EmailTemplate, qp *query.Params) (*EmailTemplate, *okta.Response, error) {
	url := "/api/v1/templates/emails"
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}

	var temp *EmailTemplate
	resp, err := m.RequestExecutor.Do(ctx, req, &temp)
	if err != nil {
		return nil, resp, err
	}
	return temp, resp, err
}

func (m *APISupplement) UpdateEmailTemplate(ctx context.Context, id string, body EmailTemplate, qp *query.Params) (*EmailTemplate, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/templates/emails/%s", id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}

	var temp *EmailTemplate
	resp, err := m.RequestExecutor.Do(ctx, req, &temp)
	if err != nil {
		return nil, resp, err
	}
	return temp, resp, err
}

func (m *APISupplement) GetEmailTemplate(ctx context.Context, id string) (*EmailTemplate, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/templates/emails/%s", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var temp *EmailTemplate
	resp, err := m.RequestExecutor.Do(ctx, req, &temp)
	if err != nil {
		return nil, resp, err
	}
	return temp, resp, err
}

func (m *APISupplement) DeleteEmailTemplate(ctx context.Context, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/templates/emails/%s", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}
