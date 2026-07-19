// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// EventHook represents the EventHook schema
type EventHook struct {
	Channel interface{} `json:"channel"`
	Events interface{} `json:"events"`
	// Date of the last event hook update
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Status of the event hook
	Status string `json:"status,omitempty"`
	VerificationStatus interface{} `json:"verificationStatus,omitempty"`
	// Timestamp of the event hook creation
	Created *time.Time `json:"created,omitempty"`
	// The ID of the user who created the event hook
	CreatedBy string `json:"createdBy,omitempty"`
	// Description of the event hook
	Description string `json:"description,omitempty"`
	// Unique key for the event hook
	ID string `json:"id,omitempty"`
	// Display name for the event hook
	Name string `json:"name"`
	Links interface{} `json:"_links,omitempty"`
}
