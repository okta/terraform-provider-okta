// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlRequestAlgorithm represents the SamlRequestAlgorithm schema
// Algorithm settings used to secure an `<AuthnRequest>` message
type SamlRequestAlgorithm struct {
	Digest interface{} `json:"digest,omitempty"`
	Signature interface{} `json:"signature,omitempty"`
}
