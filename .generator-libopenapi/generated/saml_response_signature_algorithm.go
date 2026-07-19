// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlResponseSignatureAlgorithm represents the SamlResponseSignatureAlgorithm schema
// XML digital Signature Algorithm settings for verifying `<SAMLResponse>` messages and `<Assertion>` elements from the IdP
type SamlResponseSignatureAlgorithm struct {
	Algorithm interface{} `json:"algorithm,omitempty"`
	Scope interface{} `json:"scope,omitempty"`
}
