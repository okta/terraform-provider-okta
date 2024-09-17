// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/okta/terraform-provider-okta/sdk/query"
)

type AuthenticatorResource resource

type Authenticator struct {
	Links       interface{}            `json:"_links,omitempty"`
	Created     *time.Time             `json:"created,omitempty"`
	Id          string                 `json:"id,omitempty"`
	Key         string                 `json:"key,omitempty"`
	LastUpdated *time.Time             `json:"lastUpdated,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Provider    *AuthenticatorProvider `json:"provider,omitempty"`
	Settings    *AuthenticatorSettings `json:"settings,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Type        string                 `json:"type,omitempty"`
}

type OTP struct {
	Settings *AuthenticatorSettingsOTP `json:"settings"`
}

func (m *AuthenticatorResource) GetAuthenticator(ctx context.Context, authenticatorId string) (*Authenticator, *Response, error) {
	url := fmt.Sprintf("/api/v1/authenticators/%v", authenticatorId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var authenticator *Authenticator

	resp, err := rq.Do(ctx, req, &authenticator)
	if err != nil {
		return nil, resp, err
	}

	return authenticator, resp, nil
}

// Updates an authenticator
func (m *AuthenticatorResource) UpdateAuthenticator(ctx context.Context, authenticatorId string, body Authenticator) (*Authenticator, *Response, error) {
	url := fmt.Sprintf("/api/v1/authenticators/%v", authenticatorId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var authenticator *Authenticator

	resp, err := rq.Do(ctx, req, &authenticator)
	if err != nil {
		return nil, resp, err
	}

	return authenticator, resp, nil
}

// List Authenticators
func (m *AuthenticatorResource) ListAuthenticators(ctx context.Context) ([]*Authenticator, *Response, error) {
	url := "/api/v1/authenticators"

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var authenticator []*Authenticator

	resp, err := rq.Do(ctx, req, &authenticator)
	if err != nil {
		return nil, resp, err
	}

	return authenticator, resp, nil
}

// Create Authenticator
func (m *AuthenticatorResource) CreateAuthenticator(ctx context.Context, body Authenticator, qp *query.Params) (*Authenticator, *Response, error) {
	url := "/api/v1/authenticators"
	if qp != nil {
		url = url + qp.String()
	}

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var authenticator *Authenticator

	resp, err := rq.Do(ctx, req, &authenticator)
	if err != nil {
		return nil, resp, err
	}

	return authenticator, resp, nil
}

func (m *AuthenticatorResource) SetSettingsOTP(ctx context.Context, body OTP, authenticatorId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/authenticators/%v/methods/otp", authenticatorId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}

	var otp *OTP

	resp, err := rq.Do(ctx, req, &otp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *AuthenticatorResource) ActivateAuthenticator(ctx context.Context, authenticatorId string) (*Authenticator, *Response, error) {
	url := fmt.Sprintf("/api/v1/authenticators/%v/lifecycle/activate", authenticatorId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var authenticator *Authenticator

	resp, err := rq.Do(ctx, req, &authenticator)
	if err != nil {
		return nil, resp, err
	}

	return authenticator, resp, nil
}

func (m *AuthenticatorResource) DeactivateAuthenticator(ctx context.Context, authenticatorId string) (*Authenticator, *Response, error) {
	url := fmt.Sprintf("/api/v1/authenticators/%v/lifecycle/deactivate", authenticatorId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var authenticator *Authenticator

	resp, err := rq.Do(ctx, req, &authenticator)
	if err != nil {
		return nil, resp, err
	}

	return authenticator, resp, nil
}
