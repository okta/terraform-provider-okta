// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ClientJsonEncryptionKeyRequest represents the OAuth2ClientJsonEncryptionKeyRequest schema
// A [JSON Web Key (JWK)](https://tools.ietf.org/html/rfc7517) is a JSON representation of a cryptographic key. Okta uses an encryption key to encrypt an ID token JWT minted by the org authorization s...
type OAuth2ClientJsonEncryptionKeyRequest struct {
	// Algorithm used in the key
	Alg string `json:"alg,omitempty"`
	// RSA key value (exponent) for key binding
	E string `json:"e,omitempty"`
	// Cryptographic algorithm family for the certificate's key pair
	Kty string `json:"kty,omitempty"`
	// RSA key value (modulus) for key binding
	N string `json:"n,omitempty"`
	// Acceptable use of the JSON Web Key
	Use string `json:"use,omitempty"`
}
