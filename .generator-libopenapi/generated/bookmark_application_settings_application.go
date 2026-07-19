// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// BookmarkApplicationSettingsApplication represents the BookmarkApplicationSettingsApplication schema
type BookmarkApplicationSettingsApplication struct {
	// Would you like Okta to add an integration for this app?
	RequestIntegration bool `json:"requestIntegration,omitempty"`
	// The URL of the launch page for this app
	Url string `json:"url"`
}
