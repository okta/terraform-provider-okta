// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// ResourceSetResource represents the ResourceSetResource schema
type ResourceSetResource struct {
	// Unique ID of the resource set resource object
	ID string `json:"id,omitempty"`
	// Timestamp when this object was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// The Okta Resource Name (ORN) of the resource
	Orn string `json:"orn,omitempty"`
	// Related discoverable resources
	Links interface{} `json:"_links,omitempty"`
	Conditions interface{} `json:"conditions,omitempty"`
	// Timestamp when the resource set resource object was created
	Created *time.Time `json:"created,omitempty"`
}
