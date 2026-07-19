// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IdPKeyCredential represents the IdPKeyCredential schema
// A [JSON Web Key](https://tools.ietf.org/html/rfc7517) for a signature or encryption credential for an IdP
type IdPKeyCredential struct {
	Created interface{} `json:"created,omitempty"`
	// The exponent value for the RSA public key
	E string `json:"e,omitempty"`
	// Unique identifier for the key
	Kid string `json:"kid,omitempty"`
	LastUpdated interface{} `json:"lastUpdated,omitempty"`
	// The modulus value for the RSA public key
	N string `json:"n,omitempty"`
	// Intended use of the public key
	Use string `json:"use,omitempty"`
	X5c interface{} `json:"x5c,omitempty"`
	ExpiresAt interface{} `json:"expiresAt,omitempty"`
	// Identifies the cryptographic algorithm family used with the key
	Kty string `json:"kty,omitempty"`
	// Base64url-encoded SHA-256 thumbprint of the DER encoding of an X.509 certificate
	X5t#S256 string `json:"x5t#S256,omitempty"`
}
