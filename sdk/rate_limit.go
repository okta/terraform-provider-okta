package sdk

import (
	"context"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type ClientRateLimitMode struct {
	Mode                 string                         `json:"mode"`
	GranularModeSettings *RateLimitGranularModeSettings `json:"granularModeSettings"`
}

type RateLimitGranularModeSettings struct {
	OAuth2Authorize string `json:"OAUTH2_AUTHORIZE"`
	LoginPage       string `json:"LOGIN_PAGE"`
}

type RateLimitingCommunications struct {
	RateLimitNotification *bool `json:"rateLimitNotification"`
}

func (m *APISupplement) SetClientBasedRateLimiting(ctx context.Context, body ClientRateLimitMode) (*ClientRateLimitMode, *okta.Response, error) {
	url := "/api/v1/internal/rateLimits/clientRateLimitMode"
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var rateLimitMode *ClientRateLimitMode
	resp, err := m.RequestExecutor.Do(ctx, req, &rateLimitMode)
	if err != nil {
		return nil, resp, err
	}
	return rateLimitMode, resp, nil
}

func (m *APISupplement) GetClientBasedRateLimiting(ctx context.Context) (*ClientRateLimitMode, *okta.Response, error) {
	url := "/api/v1/internal/rateLimits/clientRateLimitMode"
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var rateLimitMode *ClientRateLimitMode
	resp, err := m.RequestExecutor.Do(ctx, req, &rateLimitMode)
	if err != nil {
		return nil, resp, err
	}
	return rateLimitMode, resp, nil
}

func (m *APISupplement) SetRateLimitingCommunications(ctx context.Context, body RateLimitingCommunications) (*RateLimitingCommunications, *okta.Response, error) {
	url := "/api/internal/orgSettings/rateLimitNotificationSetting"
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var communications *RateLimitingCommunications
	resp, err := m.RequestExecutor.Do(ctx, req, &communications)
	if err != nil {
		return nil, resp, err
	}
	return communications, resp, nil
}

func (m *APISupplement) GetRateLimitingCommunications(ctx context.Context) (*RateLimitingCommunications, *okta.Response, error) {
	url := "/api/internal/orgSettings/rateLimitNotificationSetting"
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var communications *RateLimitingCommunications
	resp, err := m.RequestExecutor.Do(ctx, req, &communications)
	if err != nil {
		return nil, resp, err
	}
	return communications, resp, nil
}
