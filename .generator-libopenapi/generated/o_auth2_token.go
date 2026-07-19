// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// OAuth2Token represents the OAuth2Token schema
type OAuth2Token struct {
	// ID of the Token object
	ID string `json:"id,omitempty"`
	// Name of scopes attached to the Token
	Scopes []string `json:"scopes,omitempty"`
	UserId string `json:"userId,omitempty"`
	// Embedded resources related to the object if the `expand` query parameter is specified
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
	// Client ID
	ClientId string `json:"clientId,omitempty"`
	Created interface{} `json:"created,omitempty"`
	// Expiration time of the OAuth 2.0 Token
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	// The complete URL of the authorization server that issued the Token
	Issuer string `json:"issuer,omitempty"`
	LastUpdated interface{} `json:"lastUpdated,omitempty"`
	Status interface{} `json:"status,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
