package sdk

import (
	"context"
	"fmt"
	"net/http"
)

type HotpFactorProfile struct {
	ID       string                    `json:"id"`
	Default  bool                      `json:"default"`
	Name     string                    `json:"name"`
	Settings HotpFactorProfileSettings `json:"settings"`
}

type HotpFactorProfileSettings struct {
	TimeBased                   bool   `json:"timeBased"`
	OtpLength                   int    `json:"otpLength"`
	TimeStep                    int    `json:"timeStep"`
	AcceptableAdjacentIntervals int    `json:"acceptableAdjacentIntervals"`
	Encoding                    string `json:"encoding"`
	HmacAlgorithm               string `json:"hmacAlgorithm"`
}

func (m *APISupplement) GetHotpFactorProfile(ctx context.Context, id string) (*HotpFactorProfile, *Response, error) {
	url := fmt.Sprintf("/api/v1/org/factors/hotp/profiles/%v", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var profile *HotpFactorProfile
	resp, err := m.RequestExecutor.Do(ctx, req, &profile)
	if err != nil {
		return nil, resp, err
	}
	return profile, resp, nil
}

func (m *APISupplement) CreateHotpFactorProfile(ctx context.Context, body HotpFactorProfile) (*HotpFactorProfile, *Response, error) {
	url := "/api/v1/org/factors/hotp/profiles"
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var profile *HotpFactorProfile
	resp, err := m.RequestExecutor.Do(ctx, req, &profile)
	if err != nil {
		return nil, resp, err
	}
	return profile, resp, nil
}

func (m *APISupplement) UpdateHotpFactorProfile(ctx context.Context, id string, body HotpFactorProfile) (*HotpFactorProfile, *Response, error) {
	url := fmt.Sprintf("/api/v1/org/factors/hotp/profiles/%v", id)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}
	var profile *HotpFactorProfile
	resp, err := m.RequestExecutor.Do(ctx, req, &profile)
	if err != nil {
		return nil, resp, err
	}
	return profile, resp, nil
}

func (m *APISupplement) DeleteHotpFactorProfile(ctx context.Context, id string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/org/factors/hotp/profiles/%v", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}
