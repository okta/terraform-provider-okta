// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Org2OrgApplicationSettingsApplication represents the Org2OrgApplicationSettingsApplication schema
// Org2Org app instance properties
type Org2OrgApplicationSettingsApplication struct {
	// The base URL of the target Okta org (for `SAML_2_0` sign-on mode)
	BaseUrl string `json:"baseUrl"`
	// Used to track and manage the state of the app's creation or the provisioning process between two Okta orgs
	CreationState string `json:"creationState,omitempty"`
	// Indicates that you don't want to use an email address as the username
	PreferUsernameOverEmail bool `json:"preferUsernameOverEmail,omitempty"`
	// An API token from the target org that's used to secure the connection between the orgs
	Token string `json:"token,omitempty"`
	// Encrypted token to enhance security
	TokenEncrypted string `json:"tokenEncrypted,omitempty"`
	// The Assertion Consumer Service (ACS) URL of the source org (for `SAML_2_0` sign-on mode)
	AcsUrl string `json:"acsUrl,omitempty"`
	// The entity ID of the SP (for `SAML_2_0` sign-on mode)
	AudRestriction string `json:"audRestriction,omitempty"`
}
