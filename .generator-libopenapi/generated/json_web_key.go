// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// JsonWebKey represents the JsonWebKey schema
type JsonWebKey struct {
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// RSA modulus value that is used by both the public and private keys and provides a link between them
	N string `json:"n,omitempty"`
	// Acceptable use of the certificate. Valid value: `sig`
	Use string `json:"use,omitempty"`
	// X.509 certificate chain that contains a chain of one or more certificates
	X5c []string `json:"x5c,omitempty"`
	// X.509 certificate SHA-256 thumbprint, which is the base64url-encoded SHA-256 thumbprint (digest) of the DER encoding of an X.509 certificate
	X5t#S256 string `json:"x5t#S256,omitempty"`
	Created interface{} `json:"created,omitempty"`
	// Unique identifier for the certificate
	Kid string `json:"kid,omitempty"`
	// Cryptographic algorithm family for the certificate's keypair. Valid value: `RSA`
	Kty string `json:"kty,omitempty"`
	// RSA key value (public exponent) for Key binding
	E string `json:"e,omitempty"`
	// Timestamp when the certificate expires
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}
