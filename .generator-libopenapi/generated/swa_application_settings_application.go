// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SwaApplicationSettingsApplication represents the SwaApplicationSettingsApplication schema
type SwaApplicationSettingsApplication struct {
	// CSS selector for the **Sign-In** button in the sign-in form (for SWA apps with the `template_swa` app name definition)
	ButtonField string `json:"buttonField"`
	// Enter the CSS selector for the extra field (for three-field SWA apps with the `template_swa3field` app name definition).
	ExtraFieldSelector string `json:"extraFieldSelector,omitempty"`
	// Enter the value for the extra field in the form (for three-field SWA apps with the `template_swa3field` app name definition).
	ExtraFieldValue string `json:"extraFieldValue,omitempty"`
	// A regular expression that further restricts targetURL to the specified regular expression
	LoginUrlRegex string `json:"loginUrlRegex,omitempty"`
	// CSS selector for the **Password** field in the sign-in form (for three-field SWA apps with the `template_swa3field` app name definition)
	PasswordSelector string `json:"passwordSelector,omitempty"`
	// The URL of the sign-in page for this app (for three-field SWA apps with the `template_swa3field` app name definition)
	TargetURL string `json:"targetURL,omitempty"`
	// CSS selector for the **Username** field in the sign-in form (for SWA apps with the `template_swa` app name definition)
	UsernameField string `json:"usernameField"`
	// CSS selector for the **Username** field in the sign-in form (for three-field SWA apps with the `template_swa3field` app name definition)
	UserNameSelector string `json:"userNameSelector,omitempty"`
	// CSS selector for the **Sign-In**  button in the sign-in form (for three-field SWA apps with the `template_swa3field` app name definition)
	ButtonSelector string `json:"buttonSelector,omitempty"`
	// CSS selector for the **Password** field in the sign-in form (for SWA apps with the `template_swa` app name definition)
	PasswordField string `json:"passwordField"`
	// The URL of the sign-in page for this app (for SWA apps with the `template_swa` app name definition)
	Url string `json:"url"`
}
