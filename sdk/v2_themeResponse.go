package sdk

type ThemeResponse struct {
	Links                             interface{} `json:"_links,omitempty"`
	BackgroundImage                   string      `json:"backgroundImage,omitempty"`
	EmailTemplateTouchPointVariant    string      `json:"emailTemplateTouchPointVariant,omitempty"`
	EndUserDashboardTouchPointVariant string      `json:"endUserDashboardTouchPointVariant,omitempty"`
	ErrorPageTouchPointVariant        string      `json:"errorPageTouchPointVariant,omitempty"`
	Favicon                           string      `json:"favicon,omitempty"`
	Id                                string      `json:"id,omitempty"`
	Logo                              string      `json:"logo,omitempty"`
	PrimaryColorContrastHex           string      `json:"primaryColorContrastHex,omitempty"`
	PrimaryColorHex                   string      `json:"primaryColorHex,omitempty"`
	SecondaryColorContrastHex         string      `json:"secondaryColorContrastHex,omitempty"`
	SecondaryColorHex                 string      `json:"secondaryColorHex,omitempty"`
	SignInPageTouchPointVariant       string      `json:"signInPageTouchPointVariant,omitempty"`
}
