// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// BulkUpsertRequestBody represents the BulkUpsertRequestBody schema
type BulkUpsertRequestBody struct {
	// The type of data to upsert into the session. Currently, only `USERS` is supported.
	EntityType string `json:"entityType,omitempty"`
	// Array of user profiles to be uploaded
	Profiles []map[string]interface{} `json:"profiles,omitempty"`
}
