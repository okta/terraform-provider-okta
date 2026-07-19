// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthorizationServerJsonWebKey represents the AuthorizationServerJsonWebKey schema
type AuthorizationServerJsonWebKey struct {
	// RSA modulus value that is used by both the public and private keys and provides a link between them
	N string `json:"n,omitempty"`
	// An `ACTIVE` Key is used to sign tokens issued by the authorization server. Supported values: `ACTIVE`, `NEXT`, or `EXPIRED`<br> A `NEXT` Key is the next Key that the authorization server uses to si...
	Status string `json:"status,omitempty"`
	// Acceptable use of the key. Valid value: `sig`
	Use string `json:"use,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// The algorithm used with the Key. Valid value: `RS256`
	Alg string `json:"alg,omitempty"`
	// RSA key value (public exponent) for Key binding
	E string `json:"e,omitempty"`
	// Unique identifier for the key
	Kid string `json:"kid,omitempty"`
	// Cryptographic algorithm family for the certificate's keypair. Valid value: `RSA`
	Kty string `json:"kty,omitempty"`
}
