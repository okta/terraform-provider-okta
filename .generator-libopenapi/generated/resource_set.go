// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// ResourceSet represents the ResourceSet schema
type ResourceSet struct {
	// Timestamp when the role was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	// Timestamp when the role was created
	Created *time.Time `json:"created,omitempty"`
	// Description of the resource set
	Description string `json:"description,omitempty"`
	// Unique ID for the resource set object
	ID string `json:"id,omitempty"`
	// Unique label for the resource set
	Label string `json:"label,omitempty"`
}
