// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OAuth2Client represents the OAuth2Client schema
type OAuth2Client struct {
	Links interface{} `json:"_links,omitempty"`
	// Unique key for the client application. The `client_id` is immutable.
	ClientId string `json:"client_id,omitempty"`
	// Human-readable string name of the client application
	ClientName string `json:"client_name,omitempty"`
	ClientUri string `json:"client_uri,omitempty"`
	// URL string that references a logo for the client consent dialog (not the sign-in dialog)
	LogoUri string `json:"logo_uri,omitempty"`
}
