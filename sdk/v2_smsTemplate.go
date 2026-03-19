// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/okta/terraform-provider-okta/sdk/query"
)

type SmsTemplateResource resource

type SmsTemplate struct {
	Created      *time.Time               `json:"created,omitempty"`
	Id           string                   `json:"id,omitempty"`
	LastUpdated  *time.Time               `json:"lastUpdated,omitempty"`
	Name         string                   `json:"name,omitempty"`
	Template     string                   `json:"template,omitempty"`
	Translations *SmsTemplateTranslations `json:"translations,omitempty"`
	Type         string                   `json:"type,omitempty"`
}

// Adds a new custom SMS template to your organization.
func (m *SmsTemplateResource) CreateSmsTemplate(ctx context.Context, body SmsTemplate) (*SmsTemplate, *Response, error) {
	url := "/api/v1/templates/sms"

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var smsTemplate *SmsTemplate

	resp, err := rq.Do(ctx, req, &smsTemplate)
	if err != nil {
		return nil, resp, err
	}

	return smsTemplate, resp, nil
}

// Fetches a specific template by &#x60;id&#x60;
func (m *SmsTemplateResource) GetSmsTemplate(ctx context.Context, templateId string) (*SmsTemplate, *Response, error) {
	url := fmt.Sprintf("/api/v1/templates/sms/%v", templateId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var smsTemplate *SmsTemplate

	resp, err := rq.Do(ctx, req, &smsTemplate)
	if err != nil {
		return nil, resp, err
	}

	return smsTemplate, resp, nil
}

// Updates the SMS template.
func (m *SmsTemplateResource) UpdateSmsTemplate(ctx context.Context, templateId string, body SmsTemplate) (*SmsTemplate, *Response, error) {
	url := fmt.Sprintf("/api/v1/templates/sms/%v", templateId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var smsTemplate *SmsTemplate

	resp, err := rq.Do(ctx, req, &smsTemplate)
	if err != nil {
		return nil, resp, err
	}

	return smsTemplate, resp, nil
}

// Removes an SMS template.
func (m *SmsTemplateResource) DeleteSmsTemplate(ctx context.Context, templateId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/templates/sms/%v", templateId)

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

// Enumerates custom SMS templates in your organization. A subset of templates can be returned that match a template type.
func (m *SmsTemplateResource) ListSmsTemplates(ctx context.Context, qp *query.Params) ([]*SmsTemplate, *Response, error) {
	url := "/api/v1/templates/sms"
	if qp != nil {
		url = url + qp.String()
	}

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var smsTemplate []*SmsTemplate

	resp, err := rq.Do(ctx, req, &smsTemplate)
	if err != nil {
		return nil, resp, err
	}

	return smsTemplate, resp, nil
}

// Updates only some of the SMS template properties:
func (m *SmsTemplateResource) PartialUpdateSmsTemplate(ctx context.Context, templateId string, body SmsTemplate) (*SmsTemplate, *Response, error) {
	url := fmt.Sprintf("/api/v1/templates/sms/%v", templateId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var smsTemplate *SmsTemplate

	resp, err := rq.Do(ctx, req, &smsTemplate)
	if err != nil {
		return nil, resp, err
	}

	return smsTemplate, resp, nil
}
