package sdk

import (
	"context"
	"fmt"
)

type EmailTemplateResource resource

type EmailTemplate struct {
	Links interface{} `json:"_links,omitempty"`
	Name  string      `json:"name,omitempty"`
}

// Fetch an email template by templateName
func (m *EmailTemplateResource) GetEmailTemplate(ctx context.Context, brandId, templateName string) (*EmailTemplate, *Response, error) {
	url := fmt.Sprintf("/api/v1/brands/%v/templates/email/%v", brandId, templateName)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var emailTemplate *EmailTemplate

	resp, err := rq.Do(ctx, req, &emailTemplate)
	if err != nil {
		return nil, resp, err
	}

	return emailTemplate, resp, nil
}
