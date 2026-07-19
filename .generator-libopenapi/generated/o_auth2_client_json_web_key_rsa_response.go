// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ClientJsonWebKeyRsaResponse represents the OAuth2ClientJsonWebKeyRsaResponse schema
// An RSA signing key
type OAuth2ClientJsonWebKeyRsaResponse struct {
	// RSA key value (exponent) for key binding
	E string `json:"e,omitempty"`
	// Cryptographic algorithm family for the certificate's key pair
	Kty string `json:"kty,omitempty"`
	// RSA key value (modulus) for key binding
	N string `json:"n,omitempty"`
}
