// Code generated from OpenAPI spec. DO NOT EDIT.
package models

import "time"

// UISchemasResponseObject represents the UISchemasResponseObject schema
type UISchemasResponseObject struct {
	// Timestamp when the UI Schema was created (ISO 86001)
	Created *time.Time `json:"created"`
	// Unique identifier for the UI Schema
	ID string `json:"id"`
	// Timestamp when the UI Schema was last modified (ISO 86001)
	LastUpdated *time.Time `json:"lastUpdated"`
	UiSchema interface{} `json:"uiSchema"`
	Links interface{} `json:"_links"`
}
