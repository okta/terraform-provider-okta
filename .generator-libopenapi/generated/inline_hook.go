// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// InlineHook represents the InlineHook schema
// An inline hook object that specifies the details of the inline hook
type InlineHook struct {
	Channel interface{} `json:"channel,omitempty"`
	// The display name of the inline hook
	Name string `json:"name,omitempty"`
	Status interface{} `json:"status,omitempty"`
	// Date of the inline hook creation
	Created *time.Time `json:"created,omitempty"`
	// The unique identifier for the inline hook
	ID string `json:"id,omitempty"`
	// Date of the last inline hook update
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Type interface{} `json:"type,omitempty"`
	// Version of the inline hook type. The currently supported version is `1.0.0`.
	Version string `json:"version,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
