// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OktaPersonalAdminFeatureSettings represents the OktaPersonalAdminFeatureSettings schema
// Defines a list of Okta Personal settings that can be enabled or disabled for the org
type OktaPersonalAdminFeatureSettings struct {
	// Allow entry points for an Okta Personal account in a Workforce org
	EnableEnduserEntryPoints bool `json:"enableEnduserEntryPoints,omitempty"`
	// Allow users to migrate apps from a Workforce account to an Okta Personal account
	EnableExportApps bool `json:"enableExportApps,omitempty"`
}
