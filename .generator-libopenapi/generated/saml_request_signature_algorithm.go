// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlRequestSignatureAlgorithm represents the SamlRequestSignatureAlgorithm schema
// XML digital Signature Algorithm settings for signing `<AuthnRequest>` messages sent to the IdP > **Note:**  The `algorithm` property is ignored when you disable request signatures (`scope` set as `...
type SamlRequestSignatureAlgorithm struct {
	Algorithm interface{} `json:"algorithm,omitempty"`
	Scope interface{} `json:"scope,omitempty"`
}
