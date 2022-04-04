package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type AppOauthGroupClaim struct {
	ValueType       string `json:"valueType,omitempty"`
	GroupFilterType string `json:"groupFilterType,omitempty"`
	Issuer          string `json:"issuer,omitempty"`
	OrgURL          string `json:"orgUrl,omitempty"`
	Audience        string `json:"audience,omitempty"`
	IssuerMode      string `json:"issuerMode,omitempty"`
	Name            string `json:"name,omitempty"`
	Value           string `json:"value,omitempty"`
}

// UpdateAppOauthGroupsClaim updated OAuth app group claim
func (m *APISupplement) UpdateAppOauthGroupsClaim(ctx context.Context, appID string, gc *AppOauthGroupClaim) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/internal/apps/%s/settings/oauth/idToken", appID)
	req, err := m.RequestExecutor.NewRequest(http.MethodPost, url, gc)
	if err != nil {
		return nil, err
	}
	return m.RequestExecutor.Do(ctx, req, nil)
}

// GetAppOauthGroupsClaim gets OAuth app group claim
func (m *APISupplement) GetAppOauthGroupsClaim(ctx context.Context, appID string) (*AppOauthGroupClaim, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/internal/apps/%s/settings/oauth/idToken", appID)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var gc *AppOauthGroupClaim
	resp, err := m.RequestExecutor.Do(ctx, req, &gc)
	if err != nil {
		return nil, resp, err
	}
	return gc, resp, err
}
