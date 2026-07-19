// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GroupProfileResult represents the GroupProfileResult schema
type GroupProfileResult struct {
	// The ID of the group
	ID string `json:"id,omitempty"`
	// Map of requested attributes and their values
	Profile map[string]interface{} `json:"profile,omitempty"`
}
