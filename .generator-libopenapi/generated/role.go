// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// Role represents the Role schema
type Role struct {
	Links interface{} `json:"_links,omitempty"`
	AssignmentType interface{} `json:"assignmentType,omitempty"`
	Description string `json:"description,omitempty"`
	ID string `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Type interface{} `json:"type,omitempty"`
	Embedded map[string]interface{} `json:"_embedded,omitempty"`
	Created *time.Time `json:"created,omitempty"`
	Status interface{} `json:"status,omitempty"`
}
