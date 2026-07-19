// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// Device represents the Device schema
type Device struct {
	Links interface{} `json:"_links,omitempty"`
	// Unique key for the device
	ID string `json:"id,omitempty"`
	Profile interface{} `json:"profile,omitempty"`
	ResourceAlternateId string `json:"resourceAlternateId,omitempty"`
	// Timestamp when the device was created
	Created *time.Time `json:"created,omitempty"`
	// Timestamp when the device record was last updated. Updates occur when Okta collects and saves device signals during authentication, and when the lifecycle state of the device changes.
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	ResourceDisplayName interface{} `json:"resourceDisplayName,omitempty"`
	// Alternate key for the `id`
	ResourceId string `json:"resourceId,omitempty"`
	ResourceType string `json:"resourceType,omitempty"`
	Status interface{} `json:"status,omitempty"`
}
