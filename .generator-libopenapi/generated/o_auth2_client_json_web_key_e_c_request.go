// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ClientJsonWebKeyECRequest represents the OAuth2ClientJsonWebKeyECRequest schema
// An EC signing key
type OAuth2ClientJsonWebKeyECRequest struct {
	// The public x coordinate for the elliptic curve point
	X string `json:"x,omitempty"`
	// The public y coordinate for the elliptic curve point
	Y string `json:"y,omitempty"`
	// Cryptographic algorithm family for the certificate's key pair
	Kty string `json:"kty,omitempty"`
}
