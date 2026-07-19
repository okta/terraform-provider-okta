// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// Permission represents the Permission schema
type Permission struct {
	// The assigned Okta [permission](/openapi/okta-management/guides/permissions)
	Label string `json:"label,omitempty"`
	// Timestamp when the permission was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	Conditions interface{} `json:"conditions,omitempty"`
	// Timestamp when the permission was assigned
	Created *time.Time `json:"created,omitempty"`
}
