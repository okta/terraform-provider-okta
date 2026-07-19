// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// CustomRole represents the CustomRole schema
type CustomRole struct {
	// Timestamp when the object was created
	Created *time.Time `json:"created,omitempty"`
	// Binding object ID
	ID string `json:"id,omitempty"`
	// Timestamp when the object was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Status interface{} `json:"status,omitempty"`
	Links interface{} `json:"_links,omitempty"`
	AssignmentType interface{} `json:"assignmentType,omitempty"`
	// Label for the role assignment
	Label string `json:"label,omitempty"`
	// Resource set ID
	Resource-set string `json:"resource-set,omitempty"`
	// Role ID
	Role string `json:"role,omitempty"`
	Type interface{} `json:"type"`
}
