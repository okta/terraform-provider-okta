// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// UserClassification represents the UserClassification schema
type UserClassification struct {
	// The timestamp when the user classification was last updated
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	Type interface{} `json:"type,omitempty"`
}
