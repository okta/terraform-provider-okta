// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PasswordCredentialHook represents the PasswordCredentialHook schema
// Specify a [password import inline hook](/openapi/okta-management/management/tags/inlinehook/createpasswordimportinlinehook) to trigger verification of the user's password the first time the user si...
type PasswordCredentialHook struct {
	// The type of password inline hook. Currently, must be set to default.
	Type string `json:"type,omitempty"`
}
