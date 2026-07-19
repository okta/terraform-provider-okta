// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ResourceServerJsonWebKey represents the OAuth2ResourceServerJsonWebKey schema
type OAuth2ResourceServerJsonWebKey struct {
	Links interface{} `json:"_links,omitempty"`
	// Timestamp when the JSON Web Key was created
	Created string `json:"created,omitempty"`
	// RSA key value (exponent) for key binding
	E string `json:"e,omitempty"`
	// The unique ID of the JSON Web Key
	ID string `json:"id,omitempty"`
	// RSA key value (modulus) for key binding
	N string `json:"n,omitempty"`
	// The status of the encryption key. You can use only an `ACTIVE` key to encrypt tokens issued by the authorization server.
	Status string `json:"status,omitempty"`
	// Unique identifier of the JSON Web Key in the Custom Authorization Server's Public JWKS
	Kid string `json:"kid,omitempty"`
	// Cryptographic algorithm family for the certificate's key pair
	Kty string `json:"kty,omitempty"`
	// Timestamp when the JSON Web Key was updated
	LastUpdated string `json:"lastUpdated,omitempty"`
	// Acceptable use of the JSON Web Key
	Use string `json:"use,omitempty"`
}
