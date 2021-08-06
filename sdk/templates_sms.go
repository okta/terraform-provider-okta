package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type (
	SmsTemplate struct {
		Id           string            `json:"id,omitempty"`
		Name         string            `json:"name,omitempty"`
		Type         string            `json:"type,omitempty"`
		Template     string            `json:"template,omitempty"`
		Created      string            `json:"created,omitempty"`
		LastUpdated  string            `json:"lastUpdated,omitempty"`
		Translations map[string]string `json:"translations,omitempty"`
	}
)

func (m *APISupplement) GetSmsTemplate(ctx context.Context, id string) (*SmsTemplate, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/templates/sms/%s", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	temp := &SmsTemplate{}
	resp, err := m.RequestExecutor.Do(ctx, req, temp)

	return temp, resp, err
}

func (m *APISupplement) CreateSmsTemplate(ctx context.Context, body SmsTemplate, qp *query.Params) (*SmsTemplate, *okta.Response, error) {
	url := "/api/v1/templates/sms"
	if qp != nil {
		url += qp.String()
	}

	req, err := m.RequestExecutor.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	temp := &SmsTemplate{}
	resp, err := m.RequestExecutor.Do(ctx, req, temp)
	return temp, resp, err
}

func (m *APISupplement) UpdateSmsTemplate(ctx context.Context, id string, body SmsTemplate, qp *query.Params) (*SmsTemplate, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/templates/sms/%s", id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}
	temp := &SmsTemplate{}
	resp, err := m.RequestExecutor.Do(ctx, req, temp)
	return temp, resp, err
}
