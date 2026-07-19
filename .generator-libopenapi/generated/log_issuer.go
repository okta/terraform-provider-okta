// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogIssuer represents the LogIssuer schema
// Describes the issuer of the authorization server when the authentication is performed through OAuth. This is the location where well-known resources regarding the details of the authorization serve...
type LogIssuer struct {
	// Varies depending on the type of authentication. If authentication is SAML 2.0, `id` is the issuer in the SAML assertion. For social login, `id` is the issuer of the token.
	ID string `json:"id,omitempty"`
	// Information on the `issuer` and source of the SAML assertion or token
	Type string `json:"type,omitempty"`
}
