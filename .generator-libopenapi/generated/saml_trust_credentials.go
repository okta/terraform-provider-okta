// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlTrustCredentials represents the SamlTrustCredentials schema
// Federation Trust Credentials for verifying assertions from the IdP
type SamlTrustCredentials struct {
	// Additional IdP key credential reference to the Okta X.509 signature certificate
	AdditionalKids []interface{} `json:"additionalKids,omitempty"`
	// URI that identifies the target Okta IdP instance (SP) for an `<Assertion>`
	Audience string `json:"audience,omitempty"`
	// URI that identifies the issuer (IdP) of a `<SAMLResponse>` message `<Assertion>` element
	Issuer string `json:"issuer,omitempty"`
	Kid interface{} `json:"kid,omitempty"`
}
