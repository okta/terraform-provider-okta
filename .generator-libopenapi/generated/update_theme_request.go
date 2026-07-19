// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// UpdateThemeRequest represents the UpdateThemeRequest schema
type UpdateThemeRequest struct {
	// Primary color contrast hex code
	PrimaryColorContrastHex string `json:"primaryColorContrastHex,omitempty"`
	// Secondary color contrast hex code
	SecondaryColorContrastHex string `json:"secondaryColorContrastHex,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	EmailTemplateTouchPointVariant interface{} `json:"emailTemplateTouchPointVariant"`
	EndUserDashboardTouchPointVariant interface{} `json:"endUserDashboardTouchPointVariant"`
	LoadingPageTouchPointVariant interface{} `json:"loadingPageTouchPointVariant,omitempty"`
	// Primary color hex code
	PrimaryColorHex string `json:"primaryColorHex"`
	// Secondary color hex code
	SecondaryColorHex string `json:"secondaryColorHex"`
	SignInPageTouchPointVariant interface{} `json:"signInPageTouchPointVariant"`
	ErrorPageTouchPointVariant interface{} `json:"errorPageTouchPointVariant"`
}
