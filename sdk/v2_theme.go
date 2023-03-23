package sdk

import (
	"context"
	"fmt"
)

type ThemeResource resource

type Theme struct {
	Links                             interface{} `json:"_links,omitempty"`
	BackgroundImage                   string      `json:"backgroundImage,omitempty"`
	EmailTemplateTouchPointVariant    string      `json:"emailTemplateTouchPointVariant,omitempty"`
	EndUserDashboardTouchPointVariant string      `json:"endUserDashboardTouchPointVariant,omitempty"`
	ErrorPageTouchPointVariant        string      `json:"errorPageTouchPointVariant,omitempty"`
	PrimaryColorContrastHex           string      `json:"primaryColorContrastHex,omitempty"`
	PrimaryColorHex                   string      `json:"primaryColorHex,omitempty"`
	SecondaryColorContrastHex         string      `json:"secondaryColorContrastHex,omitempty"`
	SecondaryColorHex                 string      `json:"secondaryColorHex,omitempty"`
	SignInPageTouchPointVariant       string      `json:"signInPageTouchPointVariant,omitempty"`
}

// Fetches a theme for a brand
func (m *ThemeResource) GetBrandTheme(ctx context.Context, brandId string, themeId string) (*ThemeResponse, *Response, error) {
	url := fmt.Sprintf("/api/v1/brands/%v/themes/%v", brandId, themeId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var themeResponse *ThemeResponse

	resp, err := rq.Do(ctx, req, &themeResponse)
	if err != nil {
		return nil, resp, err
	}

	return themeResponse, resp, nil
}

// Updates a theme for a brand
func (m *ThemeResource) UpdateBrandTheme(ctx context.Context, brandId string, themeId string, body Theme) (*ThemeResponse, *Response, error) {
	url := fmt.Sprintf("/api/v1/brands/%v/themes/%v", brandId, themeId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var themeResponse *ThemeResponse

	resp, err := rq.Do(ctx, req, &themeResponse)
	if err != nil {
		return nil, resp, err
	}

	return themeResponse, resp, nil
}
