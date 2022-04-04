package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type Captcha struct {
	Id        string      `json:"id,omitempty"`
	Name      string      `json:"name"`
	SiteKey   string      `json:"siteKey"`
	SecretKey string      `json:"secretKey"`
	Type      string      `json:"type"`
	Links     interface{} `json:"_links,omitempty"`
}

func (m *APISupplement) CreateCaptcha(ctx context.Context, body Captcha) (*Captcha, *okta.Response, error) {
	url := "/api/v1/captchas"
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var captcha *Captcha
	resp, err := m.RequestExecutor.Do(ctx, req, &captcha)
	if err != nil {
		return nil, resp, err
	}
	return captcha, resp, nil
}

func (m *APISupplement) GetCaptcha(ctx context.Context, id string) (*Captcha, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/captchas/%s", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var captcha *Captcha
	resp, err := m.RequestExecutor.Do(ctx, req, &captcha)
	if err != nil {
		return nil, resp, err
	}
	return captcha, resp, nil
}

func (m *APISupplement) UpdateCaptcha(ctx context.Context, id string, body Captcha) (*Captcha, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/captchas/%s", id)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}
	var captcha *Captcha
	resp, err := m.RequestExecutor.Do(ctx, req, &captcha)
	if err != nil {
		return nil, resp, err
	}
	return captcha, resp, nil
}

func (m *APISupplement) DeleteCaptcha(ctx context.Context, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/captchas/%s", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}
