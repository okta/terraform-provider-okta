// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ClientJsonWebKeyResponseBase represents the OAuth2ClientJsonWebKeyResponseBase schema
type OAuth2ClientJsonWebKeyResponseBase struct {
	// Timestamp when the OAuth 2.0 client JSON Web Key was created
	Created string `json:"created,omitempty"`
	// The unique ID of the OAuth client JSON Web Key
	ID string `json:"id,omitempty"`
	// Timestamp when the OAuth 2.0 client JSON Web Key was updated
	LastUpdated string `json:"lastUpdated,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
