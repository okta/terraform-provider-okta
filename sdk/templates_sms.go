package sdk

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta/query"

	"github.com/okta/okta-sdk-golang/v2/okta"
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

func (m *ApiSupplement) ListSmsTemplates() ([]*SmsTemplate, *okta.Response, error) {
	url := "/api/v1/templates/sms"
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var temp []*SmsTemplate
	resp, err := m.RequestExecutor.Do(context.Background(), req, &temp)
	return temp, resp, err
}

func (m *ApiSupplement) GetSmsTemplate(id string) (*SmsTemplate, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/templates/sms/%s", id)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	temp := &SmsTemplate{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, temp)

	return temp, resp, err
}

func (m *ApiSupplement) CreateSmsTemplate(body SmsTemplate, qp *query.Params) (*SmsTemplate, *okta.Response, error) {
	url := "/api/v1/templates/sms"
	if qp != nil {
		url += qp.String()
	}

	req, err := m.RequestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}
	temp := &SmsTemplate{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, temp)
	return temp, resp, err
}

func (m *ApiSupplement) UpdateSmsTemplate(id string, body SmsTemplate, qp *query.Params) (*SmsTemplate, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/templates/sms/%s", id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.RequestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}
	temp := &SmsTemplate{}
	resp, err := m.RequestExecutor.Do(context.Background(), req, temp)
	return temp, resp, err
}

func (m *ApiSupplement) DeleteSmsTemplate(id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/templates/sms/%s", id)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(context.Background(), req, nil)
}
