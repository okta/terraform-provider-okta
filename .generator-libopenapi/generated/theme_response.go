// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ThemeResponse represents the ThemeResponse schema
type ThemeResponse struct {
	ID string `json:"id,omitempty"`
	LoadingPageTouchPointVariant interface{} `json:"loadingPageTouchPointVariant,omitempty"`
	SignInPageTouchPointVariant interface{} `json:"signInPageTouchPointVariant,omitempty"`
	ErrorPageTouchPointVariant interface{} `json:"errorPageTouchPointVariant,omitempty"`
	Favicon string `json:"favicon,omitempty"`
	Logo string `json:"logo,omitempty"`
	// Primary color contrast hex code
	PrimaryColorContrastHex string `json:"primaryColorContrastHex,omitempty"`
	// Primary color hex code
	PrimaryColorHex string `json:"primaryColorHex,omitempty"`
	// Secondary color contrast hex code
	SecondaryColorContrastHex string `json:"secondaryColorContrastHex,omitempty"`
	// Secondary color hex code
	SecondaryColorHex string `json:"secondaryColorHex,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	BackgroundImage string `json:"backgroundImage,omitempty"`
	EmailTemplateTouchPointVariant interface{} `json:"emailTemplateTouchPointVariant,omitempty"`
	EndUserDashboardTouchPointVariant interface{} `json:"endUserDashboardTouchPointVariant,omitempty"`
}
