// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlSettings represents the SamlSettings schema
// Advanced settings for the SAML 2.0 protocol
type SamlSettings struct {
	// Set to `true` to have Okta send a logout request to the upstream IdP when a user signs out of Okta or a downstream app.
	ParticipateSlo bool `json:"participateSlo,omitempty"`
	// Determines if the IdP should send the application context as `<OktaAppInstanceId>` and `<OktaAppName>` in the `<saml2p:Extensions>` element of the `<AuthnRequest>` message
	SendApplicationContext bool `json:"sendApplicationContext,omitempty"`
	// Determines if the IdP should persist account linking when the incoming assertion NameID format is `urn:oasis:names:tc:SAML:2.0:nameid-format:persistent`
	HonorPersistentNameId bool `json:"honorPersistentNameId,omitempty"`
	NameFormat interface{} `json:"nameFormat,omitempty"`
}
