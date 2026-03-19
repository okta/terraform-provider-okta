// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"net/http"
)

type OrgWideCaptchaSettings struct {
	CaptchaId    *string     `json:"captchaId"`
	EnabledPages []string    `json:"enabledPages"`
	Links        interface{} `json:"_links,omitempty"`
}

func (m *APISupplement) GetOrgWideCaptchaSettings(ctx context.Context) (*OrgWideCaptchaSettings, *Response, error) {
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

func (m *APISupplement) UpdateOrgWideCaptchaSettings(ctx context.Context, body OrgWideCaptchaSettings) (*OrgWideCaptchaSettings, *Response, error) {
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

func (m *APISupplement) DeleteOrgWideCaptchaSettings(ctx context.Context) (*Response, error) {
	url := "/api/v1/org/captcha"
	req, err := m.RequestExecutor.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}
