// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ClientJsonWebKeyRequestBody represents the OAuth2ClientJsonWebKeyRequestBody schema
type OAuth2ClientJsonWebKeyRequestBody struct {
	// Cryptographic algorithm family for the certificate's key pair
	Kty string `json:"kty,omitempty"`
	// RSA key value (modulus) for key binding
	N string `json:"n,omitempty"`
	// Status of the OAuth 2.0 client JSON Web Key
	Status string `json:"status,omitempty"`
	// Acceptable use of the JSON Web Key
	Use string `json:"use,omitempty"`
	// Algorithm used in the key
	Alg string `json:"alg,omitempty"`
	// RSA key value (exponent) for key binding
	E string `json:"e,omitempty"`
	// Unique identifier of the JSON Web Key in the OAuth 2.0 client's JWKS
	Kid string `json:"kid,omitempty"`
}
