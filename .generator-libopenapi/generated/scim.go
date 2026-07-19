// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Scim represents the Scim schema
// SCIM configuration details
type Scim struct {
	// The authentication mode for requests to your SCIM server  | authMode | Description | | -------- | ----------- | | `header` | Uses authorization header with a customer-provided token value in the fo...
	AuthMode string `json:"authMode"`
	// The base URL that Okta uses to send outbound calls to your SCIM server. Only the HTTPS protocol is supported. You can use the app-level variables defined in the `config` array for the base URL. For...
	BaseUri string `json:"baseUri"`
	EntitlementTypes interface{} `json:"entitlementTypes,omitempty"`
	// SCIM server schema configuration
	ScimServerConfig map[string]interface{} `json:"scimServerConfig"`
	// The URL to your customer-facing instructions for configuring your SCIM integration. See [Customer configuration document guidelines](https://developer.okta.com/docs/guides/submit-app-prereq/main/#c...
	SetupInstructionsUri string `json:"setupInstructionsUri"`
}
