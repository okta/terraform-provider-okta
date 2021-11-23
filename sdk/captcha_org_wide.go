package sdk

import (
	"context"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type OrgWideCaptchaSettings struct {
	CaptchaId    *string     `json:"captchaId"`
	EnabledPages []string    `json:"enabledPages"`
	Links        interface{} `json:"_links,omitempty"`
}

func (m *APISupplement) GetOrgWideCaptchaSettings(ctx context.Context) (*OrgWideCaptchaSettings, *okta.Response, error) {
	url := "/api/v1/org/captcha"
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var captcha *OrgWideCaptchaSettings
	resp, err := m.RequestExecutor.Do(ctx, req, &captcha)
	if err != nil {
		return nil, resp, err
	}
	return captcha, resp, nil
}

func (m *APISupplement) UpdateOrgWideCaptchaSettings(ctx context.Context, body OrgWideCaptchaSettings) (*OrgWideCaptchaSettings, *okta.Response, error) {
	url := "/api/v1/org/captcha"
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}
	var captcha *OrgWideCaptchaSettings
	resp, err := m.RequestExecutor.Do(ctx, req, &captcha)
	if err != nil {
		return nil, resp, err
	}
	return captcha, resp, nil
}

func (m *APISupplement) DeleteOrgWideCaptchaSettings(ctx context.Context) (*okta.Response, error) {
	url := "/api/v1/org/captcha"
	req, err := m.RequestExecutor.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}
