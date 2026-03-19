// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"fmt"
	"net/http"
)

type (
	EmailSender struct {
		ID                  string                     `json:"id,omitempty"`                  // computed
		Status              string                     `json:"status,omitempty"`              // computed
		PollingStartTime    int64                      `json:"pollingStartTime,omitempty"`    // computed
		FromName            string                     `json:"fromName,omitempty"`            // updatable
		FromAddress         string                     `json:"fromAddress,omitempty"`         // updatable
		ValidationSubdomain string                     `json:"validationSubdomain,omitempty"` // recreate
		DNSValidation       []EmailSenderDNSValidation `json:"dnsValidation,omitempty"`       // computed
		ValidationError     interface{}                `json:"validationError,omitempty"`     // computed
	}
	EmailSenderDNSValidation struct {
		RecordType        string `json:"recordType,omitempty"`
		Fqdn              string `json:"fqdn,omitempty"`
		VerificationValue string `json:"verificationValue,omitempty"`
	}
	EmailSenderValidation struct {
		PendingFromAddress      string                     `json:"pending_fromAddress,omitempty"`
		PendingFromName         string                     `json:"pending_fromName,omitempty"`
		PendingValidationDomain string                     `json:"pending_validationDomain,omitempty"`
		PendingID               string                     `json:"pending_id,omitempty"`
		PendingDNSValidation    []EmailSenderDNSValidation `json:"pending_dnsValidation,omitempty"`
		PendingStatus           string                     `json:"pending_status,omitempty"`
		PendingPollingStartTime interface{}                `json:"pending_pollingStartTime,omitempty"`
		PendingValidationError  interface{}                `json:"pending_validationError,omitempty"`
	}
	DisableActiveEmailSender struct {
		ActiveID string `json:"active_id"`
	}
	DisableInactiveEmailSender struct {
		PendingID string `json:"pending_id"`
	}
)

func (m *APISupplement) CreateEmailSender(ctx context.Context, body EmailSender) (*EmailSender, *Response, error) {
	url := "/api/v1/org/email/sender"
	req, err := m.RequestExecutor.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var sender *EmailSender
	resp, err := m.RequestExecutor.Do(ctx, req, &sender)
	if err != nil {
		return nil, resp, err
	}
	return sender, resp, err
}

func (m *APISupplement) UpdateEmailSender(ctx context.Context, body EmailSender) (*EmailSender, *Response, error) {
	url := "/api/v1/org/email/sender"
	req, err := m.RequestExecutor.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var sender *EmailSender
	resp, err := m.RequestExecutor.Do(ctx, req, &sender)
	if err != nil {
		return nil, resp, err
	}
	return sender, resp, err
}

func (m *APISupplement) GetEmailSender(ctx context.Context, id string) (*EmailSender, *Response, error) {
	url := fmt.Sprintf("/api/v1/org/email/sender/%s", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var sender *EmailSender
	resp, err := m.RequestExecutor.Do(ctx, req, &sender)
	if err != nil {
		return nil, resp, err
	}
	return sender, resp, err
}

func (m *APISupplement) DisableVerifiedEmailSender(ctx context.Context, body DisableActiveEmailSender) (*Response, error) {
	url := `/api/v1/org/email/sender/disable`
	req, err := m.RequestExecutor.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *APISupplement) DisableUnverifiedEmailSender(ctx context.Context, body DisableInactiveEmailSender) (*Response, error) {
	url := `/api/v1/org/email/sender/disable`
	req, err := m.RequestExecutor.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}

func (m *APISupplement) ValidateEmailSender(ctx context.Context, id string, body EmailSenderValidation) (*Response, error) {
	url := fmt.Sprintf("/api/v1/org/email/sender/%s/validate", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}
