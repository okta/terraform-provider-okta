// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ECKeyJWK represents the ECKeyJWK schema
// Elliptic curve key in JSON Web Key (JWK) format. It's used during enrollment to encrypt fulfillment requests to Yubico, or during activation to verify Yubico's JWS (JSON Web Signature) objects in f...
type ECKeyJWK struct {
	// The elliptic curve protocol
	Crv string `json:"crv"`
	// The unique identifier of the key
	Kid string `json:"kid"`
	// The type of public key
	Kty string `json:"kty"`
	// The intended use for the key. This value is either `enc` (encryption) during enrollment, when Okta uses the ECKeyJWK to encrypt requests to Yubico. Or it's `sig` (signature) during activation, when...
	Use string `json:"use"`
	// The public x coordinate for the elliptic curve point
	X string `json:"x"`
	// The public y coordinate for the elliptic curve point
	Y string `json:"y"`
}
