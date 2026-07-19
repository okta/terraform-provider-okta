// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ClientSecret represents the OAuth2ClientSecret schema
type OAuth2ClientSecret struct {
	// The unique ID of the OAuth 2.0 client secret
	ID string `json:"id,omitempty"`
	// Timestamp when the OAuth 2.0 client secret was updated
	LastUpdated string `json:"lastUpdated,omitempty"`
	// OAuth 2.0 client secret string hash
	SecretHash string `json:"secret_hash,omitempty"`
	// Status of the OAuth 2.0 client secret
	Status string `json:"status,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// The OAuth 2.0 client secret string
	ClientSecret string `json:"client_secret,omitempty"`
	// Timestamp when the OAuth 2.0 client secret was created
	Created string `json:"created,omitempty"`
}
