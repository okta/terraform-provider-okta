package okta

import (
	"fmt"

	"github.com/okta/okta-sdk-golang/okta"
	"github.com/okta/okta-sdk-golang/okta/query"
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

func (m *ApiSupplement) ListEmailTemplates() ([]*EmailTemplate, *okta.Response, error) {
	url := "/api/v1/templates/emails"
	req, err := m.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var auth []*EmailTemplate
	resp, err := m.requestExecutor.Do(req, &auth)
	return auth, resp, err
}

func (m *ApiSupplement) CreateEmailTemplate(id string, body EmailTemplate, qp *query.Params) (*EmailTemplate, *okta.Response, error) {
	url := "/api/v1/templates/emails"
	if qp != nil {
		url += qp.String()
	}
	req, err := m.requestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	temp := &EmailTemplate{}
	resp, err := m.requestExecutor.Do(req, temp)
	return temp, resp, err
}

func (m *ApiSupplement) UpdateEmailTemplate(id string, body EmailTemplate, qp *query.Params) (*EmailTemplate, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/templates/emails/%s", id)
	if qp != nil {
		url += qp.String()
	}
	req, err := m.requestExecutor.NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	temp := &EmailTemplate{}
	resp, err := m.requestExecutor.Do(req, temp)
	return temp, resp, err
}

func (m *ApiSupplement) GetEmailTemplate(id string) (*EmailTemplate, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/templates/emails/%s", id)
	req, err := m.requestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	temp := &EmailTemplate{}
	resp, err := m.requestExecutor.Do(req, temp)
	return temp, resp, err
}

func (m *ApiSupplement) DeleteEmailTemplate(id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/templates/emails/%s", id)
	req, err := m.requestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := m.requestExecutor.Do(req, nil)
	return resp, err
}
