// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ClientJsonWebKeyECResponse represents the OAuth2ClientJsonWebKeyECResponse schema
// An EC signing key
type OAuth2ClientJsonWebKeyECResponse struct {
	// Cryptographic algorithm family for the certificate's key pair
	Kty string `json:"kty,omitempty"`
	// The public x coordinate for the elliptic curve point
	X string `json:"x,omitempty"`
	// The public y coordinate for the elliptic curve point
	Y string `json:"y,omitempty"`
}
