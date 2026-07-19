// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// _embedded represents the _embedded schema
// The Public Key Details are defined in the `_embedded` property of the Key object.
type _embedded struct {
	// RSA key value (modulus) for key binding
	N string `json:"n,omitempty"`
	// Acceptable use of the certificate
	Use string `json:"use,omitempty"`
	// Algorithm used in the key
	Alg string `json:"alg,omitempty"`
	// RSA key value (exponent) for key binding
	E string `json:"e,omitempty"`
	// Unique identifier for the certificate
	Kid string `json:"kid,omitempty"`
	// Cryptographic algorithm family for the certificate's keypair
	Kty string `json:"kty,omitempty"`
}
