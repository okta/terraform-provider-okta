// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OrgCreationAdmin represents the OrgCreationAdmin schema
// Profile and credential information for the first super admin user of the child org. If you plan to configure and manage the org programmatically, create a system user with a dedicated email address...
type OrgCreationAdmin struct {
	// Specifies primary authentication and recovery credentials for a user. Credential types and requirements vary depending on the provider and security policy of the org.
	Credentials map[string]interface{} `json:"credentials,omitempty"`
	// Specifies the profile attributes for the first super admin user. The minimal set of required attributes are `email`, `firstName`, `lastName`, and `login`. See [profile](/openapi/okta-management/man...
	Profile map[string]interface{} `json:"profile"`
}
