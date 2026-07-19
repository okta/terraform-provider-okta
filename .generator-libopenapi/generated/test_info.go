// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// TestInfo represents the TestInfo schema
// Integration Testing Information
type TestInfo struct {
	// Details required to test the API service integration
	ApiServiceTestConfiguration map[string]interface{} `json:"apiServiceTestConfiguration,omitempty"`
	// An email for Okta to contact your company about your integration. This email isn't shared with customers.
	EscalationSupportContact string `json:"escalationSupportContact,omitempty"`
	// OIDC test details
	OidcTestConfiguration map[string]interface{} `json:"oidcTestConfiguration,omitempty"`
	// SAML test details
	SamlTestConfiguration map[string]interface{} `json:"samlTestConfiguration,omitempty"`
	// SCIM test details
	ScimTestConfiguration map[string]interface{} `json:"scimTestConfiguration,omitempty"`
	// An account on a test instance of your app with admin privileges. A test admin account is required by Okta for integration testing. During OIN QA testing, an Okta analyst uses this admin account to ...
	TestAccount map[string]interface{} `json:"testAccount,omitempty"`
}
