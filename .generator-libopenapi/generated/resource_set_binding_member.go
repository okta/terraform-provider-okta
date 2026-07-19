// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// ResourceSetBindingMember represents the ResourceSetBindingMember schema
type ResourceSetBindingMember struct {
	// Timestamp when the member was created
	Created *time.Time `json:"created,omitempty"`
	// Role resource set binding member ID
	ID string `json:"id,omitempty"`
	// Timestamp when the member was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
