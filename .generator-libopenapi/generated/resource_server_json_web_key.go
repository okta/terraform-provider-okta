// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ResourceServerJsonWebKey represents the ResourceServerJsonWebKey schema
// A [JSON Web Key (JWK)](https://tools.ietf.org/html/rfc7517) is a JSON representation of a cryptographic key. Okta can use the active key to encrypt the access token minted by the authorization serv...
type ResourceServerJsonWebKey struct {
	// The key exponent of a RSA key
	E string `json:"e,omitempty"`
	// The unique identifier of the key
	Kid string `json:"kid,omitempty"`
	Kty interface{} `json:"kty,omitempty"`
	// The modulus of the RSA key
	N string `json:"n,omitempty"`
	Status interface{} `json:"status,omitempty"`
	Use interface{} `json:"use,omitempty"`
}
