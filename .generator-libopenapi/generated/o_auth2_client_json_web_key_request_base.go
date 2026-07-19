// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ClientJsonWebKeyRequestBase represents the OAuth2ClientJsonWebKeyRequestBase schema
type OAuth2ClientJsonWebKeyRequestBase struct {
	// Unique identifier of the JSON Web Key in the OAUth 2.0 client's JWKS
	Kid string `json:"kid,omitempty"`
	// Status of the OAuth 2.0 client JSON Web Key
	Status string `json:"status,omitempty"`
}
