// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// StandardRole represents the StandardRole schema
type StandardRole struct {
	// Label for the role assignment
	Label string `json:"label,omitempty"`
	// Timestamp when the object was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Type interface{} `json:"type"`
	// Optional embedded resources for the role assignment
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
	AssignmentType interface{} `json:"assignmentType,omitempty"`
	// Timestamp when the object was created
	Created *time.Time `json:"created,omitempty"`
	// Role assignment ID
	ID string `json:"id,omitempty"`
	Status interface{} `json:"status,omitempty"`
	Links interface{} `json:"_links,omitempty"`
}
