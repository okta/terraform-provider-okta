// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// IdentitySourceSession represents the IdentitySourceSession schema
type IdentitySourceSession struct {
	// The timestamp when the identity source session was created
	Created *time.Time `json:"created,omitempty"`
	// The ID of the identity source session
	ID string `json:"id,omitempty"`
	// The ID of the custom identity source for which the session is created
	IdentitySourceId string `json:"identitySourceId,omitempty"`
	// The type of import.  All imports are `INCREMENTAL` imports.
	ImportType string `json:"importType,omitempty"`
	// The timestamp when the identity source session was created
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Status interface{} `json:"status,omitempty"`
}
