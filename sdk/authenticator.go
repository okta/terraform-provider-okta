package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type Authenticator struct {
	Type     string                 `json:"type,omitempty"`
	ID       string                 `json:"id,omitempty"`
	Key      string                 `json:"key,omitempty"`
	Status   string                 `json:"status,omitempty"`
	Name     string                 `json:"name,omitempty"`
	Settings *AuthenticatorSettings `json:"settings,omitempty"`
	Provider *AuthenticatorProvider `json:"provider,omitempty"`
}

type AuthenticatorSettings struct {
	AllowedFor             string                               `json:"allowedFor,omitempty"`
	TokenLifetimeInMinutes int                                  `json:"tokenLifetimeInMinutes,omitempty"`
	Compliance             *AuthenticatorSettingsCompliance     `json:"compliance,omitempty"`
	ChannelBinding         *AuthenticatorSettingsChannelBinding `json:"channelBinding,omitempty"`
	UserVerification       string                               `json:"userVerification,omitempty"`
	AppInstanceID          string                               `json:"appInstanceId,omitempty"`
}

type AuthenticatorSettingsCompliance struct {
	Fips string `json:"fips,omitempty"`
}

type AuthenticatorSettingsChannelBinding struct {
	Style    string `json:"style,omitempty"`
	Required string `json:"required,omitempty"`
}

type AuthenticatorProvider struct {
	Type          string                              `json:"type,omitempty"`
	Configuration *AuthenticatorProviderConfiguration `json:"configuration,omitempty"`
}

type AuthenticatorProviderConfiguration struct {
	HostName         string                                              `json:"hostName,omitempty"`
	AuthPort         int                                                 `json:"authPort,omitempty"`
	InstanceID       string                                              `json:"instanceId,omitempty"`
	SharedSecret     string                                              `json:"sharedSecret,omitempty"`
	UserNameTemplate *AuthenticatorProviderConfigurationUserNameTemplate `json:"userNameTemplate,omitempty"`
}

type AuthenticatorProviderConfigurationUserNameTemplate struct {
	Template string `json:"template,omitempty"`
}

func (m *APISupplement) ListAuthenticators(ctx context.Context) ([]*Authenticator, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authenticators")

	re := m.cloneRequestExecutor()

	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var authenticator []*Authenticator

	resp, err := re.Do(ctx, req, &authenticator)
	if err != nil {
		return nil, resp, err
	}

	return authenticator, resp, nil
}

func (m *APISupplement) GetAuthenticator(ctx context.Context, authenticatorId string) (*Authenticator, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authenticators/%v", authenticatorId)

	re := m.cloneRequestExecutor()

	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var authenticator *Authenticator

	resp, err := re.Do(ctx, req, &authenticator)
	if err != nil {
		return nil, resp, err
	}

	return authenticator, resp, nil
}

func (m *APISupplement) UpdateAuthenticator(ctx context.Context, authenticatorId string, body Authenticator) (*Authenticator, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authenticators/%v", authenticatorId)

	re := m.cloneRequestExecutor()

	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}

	var authenticator *Authenticator

	resp, err := re.Do(ctx, req, &authenticator)
	if err != nil {
		return nil, resp, err
	}

	return authenticator, resp, nil
}

func (m *APISupplement) ActivateAuthenticator(ctx context.Context, authenticatorId string) (*Authenticator, *okta.Response, error) {
	return m.lifecycleChangeAuthenticator(ctx, authenticatorId, "activate")
}

func (m *APISupplement) DeactivateAuthenticator(ctx context.Context, authenticatorId string) (*Authenticator, *okta.Response, error) {
	return m.lifecycleChangeAuthenticator(ctx, authenticatorId, "deactivate")
}

func (m *APISupplement) lifecycleChangeAuthenticator(ctx context.Context, authenticatorId, action string) (*Authenticator, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/authenticators/%v/lifecycle/%s", authenticatorId, action)
	re := m.cloneRequestExecutor()

	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var authenticator *Authenticator

	resp, err := re.Do(ctx, req, &authenticator)
	if err != nil {
		return nil, resp, err
	}

	return authenticator, resp, nil
}
