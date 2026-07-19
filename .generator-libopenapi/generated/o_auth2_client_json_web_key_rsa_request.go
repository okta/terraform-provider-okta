// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ClientJsonWebKeyRsaRequest represents the OAuth2ClientJsonWebKeyRsaRequest schema
// An RSA signing key
type OAuth2ClientJsonWebKeyRsaRequest struct {
	// RSA key value (exponent) for key binding
	E string `json:"e,omitempty"`
	// Cryptographic algorithm family for the certificate's key pair
	Kty string `json:"kty,omitempty"`
	// RSA key value (modulus) for key binding
	N string `json:"n,omitempty"`
}
