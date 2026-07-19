// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ResourceServerJsonWebKeyRequestBody represents the OAuth2ResourceServerJsonWebKeyRequestBody schema
type OAuth2ResourceServerJsonWebKeyRequestBody struct {
	// Acceptable use of the JSON Web Key
	Use string `json:"use,omitempty"`
	// RSA key value (exponent) for key binding
	E string `json:"e,omitempty"`
	// Unique identifier of the JSON web key in the custom authorization server's public JWKS
	Kid string `json:"kid,omitempty"`
	// Cryptographic algorithm family for the certificate's key pair
	Kty string `json:"kty,omitempty"`
	// RSA key value (modulus) for key binding
	N string `json:"n,omitempty"`
	// Status of the JSON Web Key
	Status string `json:"status,omitempty"`
}
