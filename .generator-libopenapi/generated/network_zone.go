// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// NetworkZone represents the NetworkZone schema
type NetworkZone struct {
	// Timestamp when the object was created
	Created *time.Time `json:"created,omitempty"`
	// Unique identifier for the Network Zone
	ID string `json:"id,omitempty"`
	// Timestamp when the object was last modified
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Unique name for this Network Zone
	Name string `json:"name"`
	Status interface{} `json:"status,omitempty"`
	// Indicates a system Network Zone: * `true` for system Network Zones * `false` for custom Network Zones  The Okta org provides the following default system Network Zones: * `LegacyIpZone` * `BlockedI...
	System bool `json:"system,omitempty"`
	Type interface{} `json:"type"`
	Usage interface{} `json:"usage,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
