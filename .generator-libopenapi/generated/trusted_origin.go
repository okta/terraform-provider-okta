// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// TrustedOrigin represents the TrustedOrigin schema
type TrustedOrigin struct {
	// Timestamp when the trusted origin was created
	Created *time.Time `json:"created,omitempty"`
	// Unique identifier for the trusted origin
	ID string `json:"id,omitempty"`
	// Timestamp when the trusted origin was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// The ID of the user who last updated the trusted origin
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	Name interface{} `json:"name,omitempty"`
	Origin interface{} `json:"origin,omitempty"`
	Status interface{} `json:"status,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// The ID of the user who created the trusted origin
	CreatedBy string `json:"createdBy,omitempty"`
	Scopes interface{} `json:"scopes,omitempty"`
}
