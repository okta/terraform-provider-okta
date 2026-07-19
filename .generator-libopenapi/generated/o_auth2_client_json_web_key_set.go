// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ClientJsonWebKeySet represents the OAuth2ClientJsonWebKeySet schema
// A JSON Web Key Set (JWKS) containing the OAuth 2.0 client's JSON Web Keys
type OAuth2ClientJsonWebKeySet struct {
	// The JSON Web Keys in this key set
	Keys []interface{} `json:"keys,omitempty"`
}
