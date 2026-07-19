// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2ClientSecretRequestBody represents the OAuth2ClientSecretRequestBody schema
type OAuth2ClientSecretRequestBody struct {
	// The OAuth 2.0 client secret string
	ClientSecret string `json:"client_secret,omitempty"`
	// Status of the OAuth 2.0 client secret
	Status string `json:"status,omitempty"`
}
