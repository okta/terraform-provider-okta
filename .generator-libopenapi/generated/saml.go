// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// Saml represents the Saml schema
// SAML configuration details
type Saml struct {
	// Defines the group attribute names for the SAML assertion statement. Okta inserts the list of Okta user groups into the attribute names in the statement.
	Groups []string `json:"groups,omitempty"`
	// List of Assertion Consumer Service (ACS) URLs. The default ACS URL is required and is indicated by a null `index` value. You can use the org-level variables you defined in the `config` array in the...
	Acs []map[string]interface{} `json:"acs"`
	// Attribute statements to appear in the Okta SAML assertion
	Claims []map[string]interface{} `json:"claims,omitempty"`
	// The URL to your customer-facing instructions for configuring your SAML integration. See [Customer configuration document guidelines](https://developer.okta.com/docs/guides/submit-app-prereq/main/#c...
	Doc string `json:"doc"`
	// Globally unique name for your SAML entity. For instance, your Identity Provider (IdP) or Service Provider (SP) URL.
	EntityId string `json:"entityId"`
}
