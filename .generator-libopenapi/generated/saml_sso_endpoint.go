// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlSsoEndpoint represents the SamlSsoEndpoint schema
// IdP's `SingleSignOnService` endpoint where Okta sends an `<AuthnRequest>` message
type SamlSsoEndpoint struct {
	Binding interface{} `json:"binding,omitempty"`
	// URI reference that indicates the address to which the `<AuthnRequest>` message is sent. The `destination` property is required if request signatures are specified. See [SAML 2.0 Request Algorithm o...
	Destination string `json:"destination,omitempty"`
	// URL of the binding-specific endpoint to send an `<AuthnRequest>` message to the IdP. The value of `url` defaults to the same value as the `sso` endpoint if omitted during creation of a new IdP inst...
	Url string `json:"url,omitempty"`
}
