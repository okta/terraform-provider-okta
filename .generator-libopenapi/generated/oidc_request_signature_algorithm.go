// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OidcRequestSignatureAlgorithm represents the OidcRequestSignatureAlgorithm schema
// Signature Algorithm settings for signing authorization requests sent to the IdP > **Note:**  The `algorithm` property is ignored when you disable request signatures (`scope` set as `NONE`).
type OidcRequestSignatureAlgorithm struct {
	Algorithm interface{} `json:"algorithm,omitempty"`
	Scope interface{} `json:"scope,omitempty"`
}
