// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// BulkDeleteRequestBody represents the BulkDeleteRequestBody schema
type BulkDeleteRequestBody struct {
	// The type of data to bulk delete in a session. Currently, only `USERS` is supported.
	EntityType string `json:"entityType,omitempty"`
	// Array of profiles to be deleted
	Profiles []interface{} `json:"profiles,omitempty"`
}
