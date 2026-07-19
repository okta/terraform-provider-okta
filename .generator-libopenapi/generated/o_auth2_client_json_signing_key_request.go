// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ClientJsonSigningKeyRequest represents the OAuth2ClientJsonSigningKeyRequest schema
// A [JSON Web Key (JWK)](https://tools.ietf.org/html/rfc7517) is a JSON representation of a cryptographic key. Okta uses signing keys to verify the signature of a JWT when provided for the `private_k...
type OAuth2ClientJsonSigningKeyRequest struct {
	// Algorithm used in the key
	Alg string `json:"alg,omitempty"`
	// Acceptable use of the JSON Web Key
	Use string `json:"use,omitempty"`
}
