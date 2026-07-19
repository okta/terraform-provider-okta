// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AutoLoginApplicationSettingsSignOn represents the AutoLoginApplicationSettingsSignOn schema
type AutoLoginApplicationSettingsSignOn struct {
	// Primary URL of the sign-in page for this app
	LoginUrl string `json:"loginUrl"`
	// Secondary URL of the sign-in page for this app
	RedirectUrl string `json:"redirectUrl,omitempty"`
}
