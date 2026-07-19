// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// OAuth2RefreshToken represents the OAuth2RefreshToken schema
type OAuth2RefreshToken struct {
	// Client ID
	ClientId string `json:"clientId,omitempty"`
	// Expiration time of the OAuth 2.0 Token
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	LastUpdated interface{} `json:"lastUpdated,omitempty"`
	// The scope names attached to the Token
	Scopes []string `json:"scopes,omitempty"`
	// The embedded resources related to the object if the `expand` query parameter is specified
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
	Created interface{} `json:"created,omitempty"`
	// ID of the Token object
	ID string `json:"id,omitempty"`
	// The complete URL of the authorization server that issued the Token
	Issuer string `json:"issuer,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// The ID of the user associated with the Token
	UserId string `json:"userId,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
